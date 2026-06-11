package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/cotishq/kronos/internal/event"
	"github.com/cotishq/kronos/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server Port")
)

type server struct {
	proto.UnimplementedEventServiceServer
	writer *kafka.Writer
	metrics *metrics
}

type metrics struct {
	eventsPublished prometheus.Counter
}

func newMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		eventsPublished: promauto.With(reg).NewCounter(prometheus.CounterOpts{
			Name: "kronos_events_published_total",
			Help: "Total number of events published to kafka",
		}),
	}
	return m
}

func(s *server) SendEvent(ctx context.Context, req *proto.SendEventRequest) (*proto.SendEventResponse, error) {
	e := event.Event{
		ID: req.GetId(),
		Type: req.GetType(),
		Payload: req.GetPayload(),
		Timestamp: time.Now(),
	}

	bytes, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	err = s.writer.WriteMessages(ctx, kafka.Message{
		Key: []byte(e.Type),
		Value: bytes,
	})

	if err != nil {
		return nil, err
	}

	log.Printf("published event: %s", e.Type)
	s.metrics.eventsPublished.Inc()
	return &proto.SendEventResponse{
		Message: "event published",
		Status: "ok",
	}, nil
}

func main() {
	flag.Parse()

	reg := prometheus.NewRegistry()
	m := newMetrics(reg)

	go func ()  {
		http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}) )
		log.Println("metrics server listening at :2112")
		http.ListenAndServe(":2112", nil)
	}()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	w := &kafka.Writer{
		Addr: kafka.TCP("localhost:9092"),
		Topic: "user-events",
		Balancer: &kafka.LeastBytes{},
	}

	defer w.Close()

	s := grpc.NewServer()
	proto.RegisterEventServiceServer(s, &server{writer: w, metrics: m})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil{
		log.Fatalf("failed to serve: %v", err)
	}
}


