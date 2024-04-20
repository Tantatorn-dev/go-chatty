package server

import (
	"context"
	"errors"
	"go-chatty/proto"
	"log"
)

var chatRooms = make(map[string][]string)

type ChatServer struct {
	proto.UnimplementedChatProtoServer
}

func (s *ChatServer) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &proto.HelloResponse{Message: "Hello " + in.GetName()}, nil
}

func (s *ChatServer) JoinRoom(ctx context.Context, in *proto.JoinRoomRequest) (*proto.JoinRoomResponse, error) {
	code := in.GetCode()

	if _, ok := chatRooms[code]; !ok {
		chatRooms[code] = []string{}
	}

	return &proto.JoinRoomResponse{Success: true}, nil
}

func (s *ChatServer) SendMessage(ctx context.Context, in *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	code := in.GetCode()
	message := in.GetMessage()

	if _, ok := chatRooms[code]; !ok {
		return &proto.SendMessageResponse{Success: false}, errors.New("room not found")
	}

	chatRooms[code] = append(chatRooms[code], message)

	return &proto.SendMessageResponse{Success: true}, nil
}

func (s *ChatServer) GetMessages(ctx context.Context, in *proto.GetMessagesRequest) (*proto.GetMessagesResponse, error) {
	code := in.GetCode()

	if _, ok := chatRooms[code]; !ok {
		return nil, errors.New("room not found")
	}

	return &proto.GetMessagesResponse{Messages: chatRooms[code]}, nil
}

func NewChatServer() *ChatServer {
	return &ChatServer{}
}
