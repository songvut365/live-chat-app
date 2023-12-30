package main

import (
	"context"
	"live-chat-app/app/internal/grpc/gen/chat/v1/chatv1connect"
	"live-chat-app/app/internal/handler"
	"live-chat-app/app/internal/manager"
	"live-chat-app/app/internal/repository"
	"live-chat-app/app/internal/service"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	address    = "localhost:50051"
	mongodbUri = "mongodb://root:1234@localhost:27017"
)

func main() {
	ctx := context.Background()

	// connect database
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbUri))
	if err != nil {
		log.Fatalf("connect mongodb error: %s", err.Error())
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("ping mongodb error: %s", err.Error())
	}
	defer mongoClient.Disconnect(ctx)

	// handler, service, manager
	chatHistoryRepository := repository.NewChatHistoryRepository(mongoClient)
	chatHistoryService := service.NewChatHistoryService(chatHistoryRepository)
	chatRoomManager := manager.NewChatRoomManager()
	chatHandler := handler.NewChatServerHandler(chatHistoryService, chatRoomManager, time.Second)

	// all path, handler from grpc server handler
	chatServicePath, chatServiceHandler := chatv1connect.NewChatServiceHandler(chatHandler)

	// routing
	mux := http.NewServeMux()
	mux.Handle(chatServicePath, chatServiceHandler)

	log.Printf("gRPC Server Listing on: %s...", address)

	// start grpc server
	http2Handler := h2c.NewHandler(mux, &http2.Server{})
	err = http.ListenAndServe(address, http2Handler)
	if err != nil {
		log.Panicf("list and service http server error: %s", err.Error())
	}
}
