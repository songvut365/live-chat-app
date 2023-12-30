package repository

import (
	"context"
	"live-chat-app/app/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	databaseName = "live_chat_app"
)

type ChatHistoryRepository interface {
	SaveMessage(ctx context.Context, chatRoom string, message model.Message) error
	GetMessages(ctx context.Context, chatRoom string) ([]model.Message, error)
}

type chatHistoryRepository struct {
	mongoClient *mongo.Client
}

func NewChatHistoryRepository(mongoClient *mongo.Client) ChatHistoryRepository {
	return &chatHistoryRepository{
		mongoClient: mongoClient,
	}
}

func (r *chatHistoryRepository) SaveMessage(ctx context.Context, chatRoom string, message model.Message) error {
	collection := r.mongoClient.Database(databaseName).Collection(chatRoom)

	_, err := collection.InsertOne(ctx, message)
	if err != nil {
		return err
	}

	return nil
}

func (r *chatHistoryRepository) GetMessages(ctx context.Context, chatRoom string) ([]model.Message, error) {
	collection := r.mongoClient.Database(databaseName).Collection(chatRoom)
	filter := bson.D{}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []model.Message = []model.Message{}
	for cursor.Next(ctx) {
		var message model.Message

		err = cursor.Decode(&message)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
