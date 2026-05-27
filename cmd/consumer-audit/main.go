package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)


func main() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic: "user-events",
		GroupID: "audit-group",
	})

	f, err := os.OpenFile("audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("failed to open audit.log:", err)
	}
	defer f.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("failed to read message", err)
			continue
		}
		
		f.WriteString(fmt.Sprintf("[audit] %s: %s\n", time.Now().Format(time.RFC3339), string(m.Key)))

	}
}
