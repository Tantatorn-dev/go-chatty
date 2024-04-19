package server

import (
	"context"
	"go-chatty/proto"
	"log"
)

type ChatServer struct {
	proto.UnimplementedChatProtoServer
}

func (s *ChatServer) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &proto.HelloResponse{Message: "Hello " + in.GetName()}, nil
}

func NewChatServer() *ChatServer {
	return &ChatServer{}
}
