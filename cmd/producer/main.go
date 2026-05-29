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
	defer w.Close()

	e := []event.Event{
		{
			ID: "evt-001",
			Type: "user_signup",
			Payload: "user123",
			Timestamp: time.Now(),
		},
		{
			ID: "evt-002",
			Type: "payment_success",
			Payload: "payment123",
			Timestamp: time.Now(),
		},
		{
			ID: "evt-003",
			Type: "file_uploaded",
			Payload: "file123",
			Timestamp: time.Now(),
		},
		
	}

	for _, evt := range e {

		bytes, err := json.Marshal(evt)
		if err != nil {
			log.Fatal("failed to marshal event:", err)
		}

		err = w.WriteMessages(context.Background(),
			kafka.Message{
				Key: []byte(evt.Type),
				Value: bytes,
			},
		)

		if err != nil {
			log.Fatal("failed to write message:", err)
			}
		
			log.Println("event published successfully:", evt.Type)
		

	}	
}