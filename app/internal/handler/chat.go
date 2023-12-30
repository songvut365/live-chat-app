package handler

import (
	"context"
	"fmt"
	chatv1 "live-chat-app/app/internal/grpc/gen/chat/v1"
	"live-chat-app/app/internal/grpc/gen/chat/v1/chatv1connect"
	"live-chat-app/app/internal/manager"
	"live-chat-app/app/internal/service"
	"log"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	ServerUserId   string = uuid.NewString()
	ServerUsername string = "server"
)

type chatServerHandler struct {
	chatHistoryService    service.ChatHistoryService
	chatRoomManager       manager.ChatRoomManager
	leaveChatCheckerDelay time.Duration
}

func NewChatServerHandler(chatHistoryService service.ChatHistoryService, chatRoomManager manager.ChatRoomManager, leaveChatCheckerDelay time.Duration) chatv1connect.ChatServiceHandler {
	return &chatServerHandler{
		chatHistoryService:    chatHistoryService,
		chatRoomManager:       chatRoomManager,
		leaveChatCheckerDelay: leaveChatCheckerDelay,
	}
}

func (h *chatServerHandler) JoinChat(ctx context.Context, req *connect.Request[chatv1.JoinChatRequest], stream *connect.ServerStream[chatv1.JoinChatResponse]) error {
	log.Printf("join chat request: %+v", req.Msg)
	roomId := manager.RoomId(req.Msg.ChatRoom.RoomId)
	userId := manager.UserId(req.Msg.User.UserId)

	h.chatRoomManager.AddConnection(manager.RoomId(roomId), userId, stream)

	_, connections := h.chatRoomManager.GetAllConnectionByUserId(userId)
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
		leaveChatUserId := <-h.chatRoomManager.LeaveChatSignal()

		if leaveChatUserId == userId {
			return nil
		}

		time.Sleep(h.leaveChatCheckerDelay)
	}
}

func (h *chatServerHandler) SendMessage(ctx context.Context, req *connect.Request[chatv1.SendMessageRequest]) (*connect.Response[chatv1.SendMessageResponse], error) {
	log.Printf("send message request: %+v", req.Msg)
	userId := manager.UserId(req.Msg.Message.Sender.UserId)

	roomId, connections := h.chatRoomManager.GetAllConnectionByUserId(userId)
	for _, connection := range connections {
		connection.Send(&chatv1.JoinChatResponse{
			Message: req.Msg.Message,
		})

	}

	err := h.chatHistoryService.SaveMessage(ctx, roomId, req.Msg.Message)
	if err != nil {
		return nil, errors.Wrap(err, "save chat history error")
	}

	return &connect.Response[chatv1.SendMessageResponse]{
		Msg: &chatv1.SendMessageResponse{
			Message: req.Msg.Message,
		},
	}, nil
}

func (h *chatServerHandler) LeaveChat(ctx context.Context, req *connect.Request[chatv1.LeaveChatRequest]) (*connect.Response[chatv1.LeaveChatResponse], error) {
	log.Printf("leave chat request: %+v", req.Msg)
	userId := manager.UserId(req.Msg.User.UserId)

	_, connections := h.chatRoomManager.GetAllConnectionByUserId(userId)
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

	h.chatRoomManager.RemoveConnection(userId)
	h.chatRoomManager.SendLeaveChatSignal(userId)

	return &connect.Response[chatv1.LeaveChatResponse]{
		Msg: &chatv1.LeaveChatResponse{
			User: req.Msg.User,
		},
	}, nil
}

func (h *chatServerHandler) GetChatHistory(ctx context.Context, req *connect.Request[chatv1.GetChatHistoryRequest]) (*connect.Response[chatv1.GetChatHistoryResponse], error) {
	chatRoom := req.Msg.ChatRoom.RoomId

	chatHistoryMessages, err := h.chatHistoryService.GetMessages(ctx, chatRoom)
	if err != nil {
		return nil, err
	}

	return &connect.Response[chatv1.GetChatHistoryResponse]{
		Msg: &chatv1.GetChatHistoryResponse{
			Message: chatHistoryMessages,
		},
	}, nil
}
