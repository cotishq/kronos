package main

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)


func main() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic: "user-events",
		GroupID: "analytics-group",
	})

	count := 0

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("failed to read message", err)
			continue
		}
		count++
		fmt.Printf("[analytics] event #%d: %s\n", count, string(m.Key))

	}
}