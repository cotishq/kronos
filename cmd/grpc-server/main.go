package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/cotishq/kronos/internal/event"
	"github.com/cotishq/kronos/proto"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server Port")
)

type server struct {
	proto.UnimplementedEventServiceServer
	writer *kafka.Writer
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
	return &proto.SendEventResponse{
		Message: "event published",
		Status: "ok",
	}, nil
}

func main() {
	flag.Parse()
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
	proto.RegisterEventServiceServer(s, &server{writer: w})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil{
		log.Fatalf("failed to serve: %v", err)
	}
}


