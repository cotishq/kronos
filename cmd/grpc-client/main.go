package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/cotishq/kronos/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewEventServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SendEvent(ctx, &proto.SendEventRequest{
		Id: "evt-001",
		Type: "user_signup",
		Payload: "user123",
	})
	if err != nil {
		log.Fatalf("could not send event: %v",err)
	}
	log.Printf("event: %s", r.GetMessage())
}
