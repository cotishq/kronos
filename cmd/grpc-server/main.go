package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/cotishq/kronos/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server Port")
)

type server struct {
	proto.UnimplementedEventServiceServer
}

func(s *server) SendEvent(ctx context.Context, req *proto.SendEventRequest) (*proto.SendEventResponse, error) {
	log.Printf("received: %v", req.GetType())
	return &proto.SendEventResponse{Message: "type of event:" + req.GetType()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterEventServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil{
		log.Fatalf("failed to serve: %v", err)
	}
}


