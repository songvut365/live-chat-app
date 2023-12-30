package service

import (
	"context"
	chatv1 "live-chat-app/app/internal/grpc/gen/chat/v1"
	"live-chat-app/app/internal/model"
	"live-chat-app/app/internal/repository"
	"time"
)

type ChatHistoryService interface {
	SaveMessage(ctx context.Context, chatRoom string, message *chatv1.Message) error
	GetMessages(ctx context.Context, chatRoom string) ([]*chatv1.Message, error)
}

type chatHistoryService struct {
	chatHistoryRepository repository.ChatHistoryRepository
}

func NewChatHistoryService(chatHistoryRepository repository.ChatHistoryRepository) ChatHistoryService {
	return &chatHistoryService{
		chatHistoryRepository: chatHistoryRepository,
	}
}

func (svc *chatHistoryService) SaveMessage(ctx context.Context, chatRoom string, message *chatv1.Message) error {
	messageTimeStamp, err := time.Parse(time.RFC3339, message.Timestamp)
	if err != nil {
		return err
	}

	err = svc.chatHistoryRepository.SaveMessage(
		ctx,
		chatRoom,
		model.Message{
			Id:        message.MessageId,
			Sender:    message.Sender.Username,
			Content:   message.Content,
			Timestamp: messageTimeStamp,
		})
	if err != nil {
		return err
	}

	return nil
}

func (svc *chatHistoryService) GetMessages(ctx context.Context, chatRoom string) ([]*chatv1.Message, error) {
	var messages []*chatv1.Message = []*chatv1.Message{}

	result, err := svc.chatHistoryRepository.GetMessages(ctx, chatRoom)
	if err != nil {
		return nil, err
	}

	for _, message := range result {
		messages = append(messages, &chatv1.Message{
			MessageId: message.Id,
			Sender: &chatv1.User{
				Username: message.Sender,
			},
			Content:   message.Content,
			Timestamp: message.Timestamp.Format(time.RFC3339),
		})
	}

	return messages, nil
}
