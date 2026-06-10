package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cotishq/kronos/internal/event"
	"github.com/segmentio/kafka-go"
)

var seenEvents sync.Map

func processEvent(e event.Event, eventCounts map[string]int) error {
	eventCounts[e.Type]++
	fmt.Printf("[analytics] %s count=%d\n", e.Type, eventCounts[e.Type])
	return nil
}

func main() {
	dlq := &kafka.Writer{
		Addr: kafka.TCP("localhost:9092"),
		Topic: "dead-letter",
		Balancer: &kafka.LeastBytes{},
	}
	defer dlq.Close()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic: "user-events",
		GroupID: "analytics-group",
	})
	defer r.Close()

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

		}
	}

	fmt.Println("analytics consumer stopped")
	
}