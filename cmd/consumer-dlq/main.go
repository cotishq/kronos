package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/cotishq/kronos/internal/event"
	"github.com/segmentio/kafka-go"
)

var (
	kafkaBroker = flag.String("kafka-broker", "localhost:9092", "Kafka broker address")
)

func main()  {
	flag.Parse()
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{*kafkaBroker},
		Topic: "dead-letter",
		GroupID: "dlq-group",
	})

	defer r.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	for {
		m, err := r.ReadMessage(ctx)
		var e event.Event
		if err != nil {
			if ctx.Err() != nil {
				fmt.Println("\nshutting down dlq consumer")
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

		fmt.Printf("[DLQ] failed event received: %s id=%s\n", e.Type, e.ID)
	}
}