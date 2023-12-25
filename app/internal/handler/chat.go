package handler

import (
	"context"
	"fmt"
	chatv1 "live-chat-app/app/internal/grpc/gen/chat/v1"
	"live-chat-app/app/internal/manager"
	"log"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
)

var (
	ServerUserId   string = uuid.NewString()
	ServerUsername string = "server"
)

type ChatServerHandler interface {
	JoinChat(context.Context, *connect.Request[chatv1.JoinChatRequest], *connect.ServerStream[chatv1.JoinChatResponse]) error
	SendMessage(context.Context, *connect.Request[chatv1.SendMessageRequest]) (*connect.Response[chatv1.SendMessageResponse], error)
	LeaveChat(context.Context, *connect.Request[chatv1.LeaveChatRequest]) (*connect.Response[chatv1.LeaveChatResponse], error)
}

type chatServerHandler struct {
	chatRoomManager       *manager.ChatRoomManager
	leaveChatCheckerDelay time.Duration
}

func NewChatServerHandler(chatRoomManager *manager.ChatRoomManager, leaveChatCheckerDelay time.Duration) ChatServerHandler {
	return chatServerHandler{
		chatRoomManager:       chatRoomManager,
		leaveChatCheckerDelay: leaveChatCheckerDelay,
	}
}

func (handler chatServerHandler) JoinChat(ctx context.Context, req *connect.Request[chatv1.JoinChatRequest], stream *connect.ServerStream[chatv1.JoinChatResponse]) error {
	log.Printf("join chat request: %+v", req.Msg)
	roomId := manager.RoomId(req.Msg.ChatRoom.RoomId)
	userId := manager.UserId(req.Msg.User.UserId)

	handler.chatRoomManager.AddConnection(manager.RoomId(roomId), manager.UserId(userId), stream)

	connections := handler.chatRoomManager.GetAllConnectionByUserId(manager.UserId(userId))
	for _, connection := range connections {
		connection.Send(&chatv1.JoinChatResponse{
			Message: &chatv1.Message{
				MessageId: uuid.NewString(),
				Sender: &chatv1.User{
					UserId:   ServerUserId,
					Username: ServerUsername,
				},
				Content:   fmt.Sprintf("%s join the chat", req.Msg.User.Username),
				Timestamp: time.Now().Format(time.RFC3339),
			},
		})
	}

	for {
		leaveChatUserId := <-handler.chatRoomManager.LeaveChatSignal

		if leaveChatUserId == manager.UserId(userId) {
			return nil
		}

		time.Sleep(handler.leaveChatCheckerDelay)
	}
}

func (handler chatServerHandler) SendMessage(ctx context.Context, req *connect.Request[chatv1.SendMessageRequest]) (*connect.Response[chatv1.SendMessageResponse], error) {
	log.Printf("send message request: %+v", req.Msg)
	userId := manager.UserId(req.Msg.Message.Sender.UserId)

	connections := handler.chatRoomManager.GetAllConnectionByUserId(manager.UserId(userId))
	for _, connection := range connections {
		connection.Send(&chatv1.JoinChatResponse{
			Message: req.Msg.Message,
		})
	}

	return &connect.Response[chatv1.SendMessageResponse]{
		Msg: &chatv1.SendMessageResponse{
			Message: req.Msg.Message,
		},
	}, nil
}

func (handler chatServerHandler) LeaveChat(ctx context.Context, req *connect.Request[chatv1.LeaveChatRequest]) (*connect.Response[chatv1.LeaveChatResponse], error) {
	log.Printf("leave chat request: %+v", req.Msg)
	userId := manager.UserId(req.Msg.User.UserId)

	connections := handler.chatRoomManager.GetAllConnectionByUserId(manager.UserId(userId))
	for _, connection := range connections {
		connection.Send(&chatv1.JoinChatResponse{
			Message: &chatv1.Message{
				MessageId: uuid.NewString(),
				Sender: &chatv1.User{
					UserId:   ServerUserId,
					Username: ServerUsername,
				},
				Content:   fmt.Sprintf("%s left the chat", req.Msg.User.Username),
				Timestamp: time.Now().Format(time.RFC3339),
			},
		})
	}

	handler.chatRoomManager.RemoveConnection(manager.UserId(userId))
	handler.chatRoomManager.LeaveChatSignal <- manager.UserId(userId)

	return &connect.Response[chatv1.LeaveChatResponse]{
		Msg: &chatv1.LeaveChatResponse{
			User: req.Msg.User,
		},
	}, nil
}
