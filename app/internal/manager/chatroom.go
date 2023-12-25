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

type ChatRoomManager struct {
	ChatRooms       map[RoomId]map[UserId]*connect.ServerStream[chatv1.JoinChatResponse]
	ChatRoomsMutex  sync.Mutex
	LeaveChatSignal chan UserId
}

func NewChatRoomManager() *ChatRoomManager {
	return &ChatRoomManager{
		ChatRooms:       make(map[RoomId]map[UserId]*connect.ServerStream[chatv1.JoinChatResponse]),
		LeaveChatSignal: make(chan UserId),
	}
}

func (manager *ChatRoomManager) AddConnection(roomId RoomId, userId UserId, stream *connect.ServerStream[chatv1.JoinChatResponse]) {
	manager.ChatRoomsMutex.Lock()
	defer manager.ChatRoomsMutex.Unlock()

	if manager.ChatRooms[roomId] == nil {
		manager.ChatRooms[roomId] = make(map[UserId]*connect.ServerStream[chatv1.JoinChatResponse])
	}

	manager.ChatRooms[roomId][userId] = stream

	log.Printf("chat rooms: %+v", manager.ChatRooms)
}

func (manager *ChatRoomManager) GetAllConnectionByUserId(userId UserId) []*connect.ServerStream[chatv1.JoinChatResponse] {
	manager.ChatRoomsMutex.Lock()
	defer manager.ChatRoomsMutex.Unlock()

	connections := []*connect.ServerStream[chatv1.JoinChatResponse]{}

	for _, chatRoom := range manager.ChatRooms {
		if _, ok := chatRoom[userId]; ok {
			for _, connection := range chatRoom {
				connections = append(connections, connection)
			}
		}

		continue
	}

	return connections
}

func (manager *ChatRoomManager) RemoveConnection(userId UserId) {
	manager.ChatRoomsMutex.Lock()
	defer manager.ChatRoomsMutex.Unlock()

	for roomId, connection := range manager.ChatRooms {
		if _, ok := connection[userId]; ok {
			delete(manager.ChatRooms[roomId], userId)
		}

		continue
	}

	log.Printf("chat rooms: %+v", manager.ChatRooms)
}
