package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)


func main() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic: "user-events",
		GroupID: "analytics-group",
	})
	defer r.Close()


	count := 0
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	for {
		m, err := r.ReadMessage(ctx)

		if err != nil {
			if ctx.Err() != nil {
				fmt.Println("shutting down analytics consumer")
				break
			}

			fmt.Println("failed to read message", err)
			continue
		}
		
		count++
		fmt.Printf("[analytics] event #%d: %s\n", count, string(m.Key))
		

	}

	fmt.Println("analytics consumer stopped")
	
	
}