package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cotishq/kronos/internal/event"
	"github.com/segmentio/kafka-go"
)


func main() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic: "user-events",
		GroupID: "audit-group",
	})
	defer r.Close()

	f, err := os.OpenFile("audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("failed to open audit.log:", err)
	}
	defer f.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	for {
		m, err := r.ReadMessage(ctx)
		var e event.Event
		if err != nil {
			if ctx.Err() != nil {
				fmt.Println("\nshutting down audit consumer")
				break
			}
			fmt.Println("failed to read message", err)
			continue
		}

		err = json.Unmarshal(m.Value, &e)
		if err != nil {
			fmt.Println("failed to unmarshal:", err)
			continue
		}
		
		fmt.Fprintf(f, "[audit] %s: %s\n", time.Now().Format(time.RFC3339), e.Type)

	}
	fmt.Println("audit-consumer stopped")
}
