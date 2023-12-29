package main

import (
	"context"
	"live-chat-app/app/internal/client"
	"live-chat-app/app/internal/grpc/gen/chat/v1/chatv1connect"
	"live-chat-app/app/internal/manager"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
)

const (
	address = "http://localhost:50051"
)

func main() {
	// new chat client
	chatIOManager := manager.NewChatIOManager()
	chatServiceClient := chatv1connect.NewChatServiceClient(http.DefaultClient, address, connect.WithGRPC())
	chatClient := client.NewChatClient(chatServiceClient)

	// init user id
	userId := uuid.NewString()

	// get username and chat room
	username, err := chatIOManager.ReadUsername()
	if err != nil {
		log.Fatalf("read username error: %s", err.Error())
	}

	chatRoom, err := chatIOManager.ReadChatRoom()
	if err != nil {
		log.Fatalf("read chat room error: %s", err.Error())
	}

	// join chat server
	ctx := context.Background()
	stream, err := chatClient.JoinChat(ctx, userId, username, chatRoom)
	if err != nil {
		log.Fatalf("join chat error: %s", err.Error())
	}

	// receive and display messages from chat server
	go func() {
		defer stream.Close()

		for {
			if stream.Receive() && stream.Msg().Message.Sender.Username != username {
				chatIOManager.DisplayMessage(stream.Msg().Message, username)
			}
		}
	}()

	// infinite loop for reading and sending messages
	for {
		message, err := chatIOManager.ReadMessage(username)
		if err != nil {
			log.Fatalf("\nread message error: %s", err.Error())
		}

		if message == "/exit" {
			break
		}

		_, err = chatClient.SendMessage(ctx, userId, username, message)
		if err != nil {
			log.Fatalf("send message error: %s", err.Error())
		}

		time.Sleep(time.Millisecond * 100)
	}

	// leave chat server
	_, err = chatClient.LeaveChat(ctx, userId, username)
	if err != nil {
		log.Fatalf("leave chat error: %s", err)
	}
}
