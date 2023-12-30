package manager

import (
	chatv1 "live-chat-app/app/internal/grpc/gen/chat/v1"
	"log"
	"sync"

	"connectrpc.com/connect"
)

type (
	RoomId string
	UserId string
)

type ChatRoomManager interface {
	AddConnection(roomId RoomId, userId UserId, stream *connect.ServerStream[chatv1.JoinChatResponse])
	GetAllConnectionByUserId(userId UserId) (chatRoomName string, connections []*connect.ServerStream[chatv1.JoinChatResponse])
	RemoveConnection(userId UserId)
	LeaveChatSignal() <-chan UserId
	SendLeaveChatSignal(userId UserId)
}

type chatRoomManager struct {
	ChatRooms       map[RoomId]map[UserId]*connect.ServerStream[chatv1.JoinChatResponse]
	leaveChatSignal chan UserId
	chatRoomsMutex  sync.Mutex
}

func NewChatRoomManager() ChatRoomManager {
	return &chatRoomManager{
		ChatRooms:       make(map[RoomId]map[UserId]*connect.ServerStream[chatv1.JoinChatResponse]),
		leaveChatSignal: make(chan UserId),
	}
}

func (manager *chatRoomManager) AddConnection(roomId RoomId, userId UserId, stream *connect.ServerStream[chatv1.JoinChatResponse]) {
	manager.chatRoomsMutex.Lock()
	defer manager.chatRoomsMutex.Unlock()

	if manager.ChatRooms[roomId] == nil {
		manager.ChatRooms[roomId] = make(map[UserId]*connect.ServerStream[chatv1.JoinChatResponse])
	}

	manager.ChatRooms[roomId][userId] = stream

	log.Printf("chat rooms: %+v", manager.ChatRooms)
}

func (manager *chatRoomManager) GetAllConnectionByUserId(userId UserId) (chatRoomName string, connections []*connect.ServerStream[chatv1.JoinChatResponse]) {
	manager.chatRoomsMutex.Lock()
	defer manager.chatRoomsMutex.Unlock()

	connections = []*connect.ServerStream[chatv1.JoinChatResponse]{}

	for roomId, chatRoom := range manager.ChatRooms {
		if _, ok := chatRoom[userId]; ok {
			chatRoomName = string(roomId)

			for _, connection := range chatRoom {
				connections = append(connections, connection)
			}
		}

		continue
	}

	return chatRoomName, connections
}

func (manager *chatRoomManager) RemoveConnection(userId UserId) {
	manager.chatRoomsMutex.Lock()
	defer manager.chatRoomsMutex.Unlock()

	for roomId, connection := range manager.ChatRooms {
		if _, ok := connection[userId]; ok {
			delete(manager.ChatRooms[roomId], userId)
		}

		continue
	}

	log.Printf("chat rooms: %+v", manager.ChatRooms)
}

func (manager *chatRoomManager) LeaveChatSignal() <-chan UserId {
	return manager.leaveChatSignal
}

func (manager *chatRoomManager) SendLeaveChatSignal(userId UserId) {
	manager.leaveChatSignal <- userId
}
