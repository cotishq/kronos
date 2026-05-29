package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/cotishq/kronos/internal/event"
	"github.com/segmentio/kafka-go"
)


func main() {
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
		eventCounts[e.Type]++
		fmt.Printf(
			"[analytics] %s count=%d\n",
			e.Type,
			eventCounts[e.Type],
			)
		

	}

	fmt.Println("analytics consumer stopped")
	
}