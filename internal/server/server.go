package server

import (
	"context"
	"fmt"
	"go-chatty/proto"

	"github.com/go-redis/redis"
)

type ChatServer struct {
	proto.UnimplementedChatProtoServer
}

func getClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return rdb
}

func (s *ChatServer) SendMessage(ctx context.Context, in *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	code := in.GetCode()
	message := in.GetMessage()

	client := getClient()

	client.XAdd(&redis.XAddArgs{
		Stream: "chat:" + code,
		Values: map[string]interface{}{"message": message},
	})

	return &proto.SendMessageResponse{Success: true}, nil
}

func (s *ChatServer) GetMessages(ctx context.Context, in *proto.GetMessagesRequest) (*proto.GetMessagesResponse, error) {
	code := in.GetCode()

	client := getClient()
	defer client.Close()

	// read latest 10 messages
	cmd := client.XRead(&redis.XReadArgs{
		Streams: []string{"chat:" + code, "0"},
		Count:   10,
		Block:   0,
	})

	var messages []string

	streams, err := cmd.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to read messages: %v", err)
	}

	for _, st := range streams {
		for _, message := range st.Messages {
			messages = append(messages, message.Values["message"].(string))
		}
	}

	return &proto.GetMessagesResponse{Messages: messages}, nil
}

func (s *ChatServer) JoinRoom(ctx context.Context, in *proto.JoinRoomRequest) (*proto.JoinRoomResponse, error) {
	code := in.GetCode()

	client := getClient()

	// create a new stream
	client.XAdd(&redis.XAddArgs{
		Stream: "chat:" + code,
		Values: map[string]interface{}{"message": "User joined"},
	})

	// read latest 10 messages
	cmd := client.XRead(&redis.XReadArgs{
		Streams: []string{"chat:" + code, "0"},
		Count:   10,
		Block:   0,
	})

	var messages []string

	streams, err := cmd.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to read messages: %v", err)
	}

	for _, st := range streams {
		for _, message := range st.Messages {
			messages = append(messages, message.Values["message"].(string))
		}
	}

	return &proto.JoinRoomResponse{Messages: messages}, nil
}

func NewChatServer() *ChatServer {
	return &ChatServer{}
}
