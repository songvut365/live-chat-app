package client

import (
	"context"
	chatv1 "live-chat-app/app/internal/grpc/gen/chat/v1"
	"live-chat-app/app/internal/grpc/gen/chat/v1/chatv1connect"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
)

type ChatClient interface {
	JoinChat(ctx context.Context, userId, username, chatroom string) (*connect.ServerStreamForClient[chatv1.JoinChatResponse], error)
	SendMessage(ctx context.Context, userId, username, message string) (*connect.Response[chatv1.SendMessageResponse], error)
	LeaveChat(ctx context.Context, userId, username string) (*connect.Response[chatv1.LeaveChatResponse], error)
	GetChatHistory(ctx context.Context, chatRoom string) (*connect.Response[chatv1.GetChatHistoryResponse], error)
}

type chatClient struct {
	service chatv1connect.ChatServiceClient
}

func NewChatClient(service chatv1connect.ChatServiceClient) ChatClient {
	return &chatClient{
		service: service,
	}
}

func (client *chatClient) JoinChat(ctx context.Context, userId, username, chatRoom string) (*connect.ServerStreamForClient[chatv1.JoinChatResponse], error) {
	request := &connect.Request[chatv1.JoinChatRequest]{
		Msg: &chatv1.JoinChatRequest{
			User: &chatv1.User{
				UserId:   userId,
				Username: username,
			},
			ChatRoom: &chatv1.ChatRoom{
				RoomId: chatRoom,
			},
		},
	}

	stream, err := client.service.JoinChat(ctx, request)
	if err != nil {
		return nil, err
	}

	return stream, nil
}

func (client *chatClient) SendMessage(ctx context.Context, userId, username, message string) (*connect.Response[chatv1.SendMessageResponse], error) {
	request := &connect.Request[chatv1.SendMessageRequest]{
		Msg: &chatv1.SendMessageRequest{
			Message: &chatv1.Message{
				MessageId: uuid.NewString(),
				Sender: &chatv1.User{
					UserId:   userId,
					Username: username,
				},
				Content:   message,
				Timestamp: time.Now().Format(time.RFC3339),
			},
		},
	}

	response, err := client.service.SendMessage(ctx, request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *chatClient) LeaveChat(ctx context.Context, userId, username string) (*connect.Response[chatv1.LeaveChatResponse], error) {
	request := &connect.Request[chatv1.LeaveChatRequest]{
		Msg: &chatv1.LeaveChatRequest{
			User: &chatv1.User{
				UserId:   userId,
				Username: username,
			},
		},
	}

	res, err := client.service.LeaveChat(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (client *chatClient) GetChatHistory(ctx context.Context, chatRoom string) (*connect.Response[chatv1.GetChatHistoryResponse], error) {
	request := &connect.Request[chatv1.GetChatHistoryRequest]{
		Msg: &chatv1.GetChatHistoryRequest{
			ChatRoom: &chatv1.ChatRoom{
				RoomId: chatRoom,
			},
		},
	}

	res, err := client.service.GetChatHistory(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}
