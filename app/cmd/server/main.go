package main

import (
	"live-chat-app/app/internal/grpc/gen/chat/v1/chatv1connect"
	"live-chat-app/app/internal/handler"
	"live-chat-app/app/internal/manager"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	address = "localhost:50051"
)

func main() {
	// handler, service, manager
	chatRoomManager := manager.NewChatRoomManager()
	chatHandler := handler.NewChatServerHandler(chatRoomManager, time.Second)

	// all path, handler from grpc server handler
	chatServicePath, chatServiceHandler := chatv1connect.NewChatServiceHandler(chatHandler)

	// routing
	mux := http.NewServeMux()
	mux.Handle(chatServicePath, chatServiceHandler)

	log.Printf("gRPC Server Listing on: %s...", address)

	// start grpc server
	http2Handler := h2c.NewHandler(mux, &http2.Server{})
	err := http.ListenAndServe(address, http2Handler)
	if err != nil {
		log.Panicf("list and service http server error: %s", err.Error())
	}
}
