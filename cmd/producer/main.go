package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/cotishq/kronos/internal/event"
	"github.com/segmentio/kafka-go"
)

func main() {
	w := &kafka.Writer{
		Addr: kafka.TCP("localhost:9092"),
		Topic: "user-events",
		Balancer: &kafka.LeastBytes{},
	}

	e := event.Event{
		ID: "evt-001",
		Type: "user_signup",
		Payload: "user123",
		Timestamp: time.Now(),
	}

	bytes, err := json.Marshal(e)
	if err != nil {
		log.Fatal("failed to marshal event:", err)
	}

	err = w.WriteMessages(context.Background(),
		kafka.Message{
			Key: []byte(e.Type),
			Value: bytes,
		},
)

if err != nil {
	log.Fatal("failed to write message:", err)
}
if err := w.Close(); err != nil {
	log.Fatal("failed to close writer:", err)
}

log.Println("event published successfully:", e.Type)

	
}