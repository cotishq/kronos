package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cotishq/kronos/internal/event"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/segmentio/kafka-go"
)

var seenEvents sync.Map

var (
	kafkaBroker = flag.String("kafka-broker", "localhost:9092", "Kafka broker address")
)

type metrics struct {
	eventsProcessed prometheus.Counter
	eventsDlq 		prometheus.Counter
}

func newMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		eventsProcessed: promauto.With(reg).NewCounter(prometheus.CounterOpts{
			Name: "kronos_events_processed_total",
			Help: "Total number of events processed",
		}),
		eventsDlq: promauto.With(reg).NewCounter(prometheus.CounterOpts{
			Name: "kronos_events_dlq_total",
			Help: "Total number of dlq events",
		}),
	}
	return m
}

func processEvent(e event.Event, eventCounts map[string]int) error {
	eventCounts[e.Type]++
	fmt.Printf("[analytics] %s count=%d\n", e.Type, eventCounts[e.Type])
	return nil
}

func main() {
	flag.Parse()
	dlq := &kafka.Writer{
		Addr: kafka.TCP(*kafkaBroker),
		Topic: "dead-letter",
		Balancer: &kafka.LeastBytes{},
	}
	defer dlq.Close()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{*kafkaBroker},
		Topic: "user-events",
		GroupID: "analytics-group",
	})
	defer r.Close()

	reg := prometheus.NewRegistry()
	metrics := newMetrics(reg)

	go func ()  {
		http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
		log.Println("metrics server listening at :2113")
		http.ListenAndServe(":2113", nil)
	}()

    eventCounts := make(map[string]int)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	for {
		m, err := r.ReadMessage(ctx)
		var e event.Event

		if err != nil {
			if ctx.Err() != nil {
				fmt.Println("shutting down analytics consumer")
				break
			}

			fmt.Println("failed to read message", err)
			continue
		}


		err = json.Unmarshal(m.Value, &e)
		if err != nil {
			fmt.Println("failed to unmarshal event:", err)
			continue
		}

		if _, seen := seenEvents.Load(e.ID); seen {
			fmt.Printf("[analytics] duplicate event skipped: %s\n", e.ID)
			continue
		}

		seenEvents.Store(e.ID, true)

		maxRetries := 3
		var processErr error
		for attempt := 0; attempt < maxRetries; attempt++ {
			processErr = processEvent(e, eventCounts)
			if processErr == nil {
				metrics.eventsProcessed.Inc()
				break
			}
			fmt.Printf("[analytics] attempt %d failed: %v, retrying...\n", attempt+1, processErr)
			time.Sleep(time.Duration(attempt+1) * time.Second)
		}
		if processErr != nil {
			fmt.Printf("[analytics] all retries exhausted for event: %s, sending to DLQ\n", e.ID)

			bytes, _ := json.Marshal(e)
			dlq.WriteMessages(context.Background(), kafka.Message{
				Key: []byte(e.ID),
				Value: bytes,
			})
			metrics.eventsDlq.Inc()

		}
	}

	fmt.Println("analytics consumer stopped")
	
}