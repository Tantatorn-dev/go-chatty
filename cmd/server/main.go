package main

import (
	"fmt"
	"go-chatty/internal/server"
	"go-chatty/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	port = 8090
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	chatServer := server.NewChatServer()
	proto.RegisterChatProtoServer(s, chatServer)

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
