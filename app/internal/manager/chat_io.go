package manager

import (
	"bufio"
	"fmt"
	chatv1 "live-chat-app/app/internal/grpc/gen/chat/v1"
	"os"
	"strings"
)

type ChatIOManager interface {
	ReadUsername() (string, error)
	ReadChatRoom() (string, error)
	ReadMessage(username string) (string, error)
	DisplayMessage(message *chatv1.Message, username string)
}

type chatIOManager struct {
	inputReader *bufio.Reader
}

func NewChatIOManager() ChatIOManager {
	return &chatIOManager{
		inputReader: bufio.NewReader(os.Stdin),
	}
}

func (manager *chatIOManager) ReadUsername() (string, error) {
	fmt.Printf("Enter your username: ")
	username, err := manager.inputReader.ReadString('\n')
	if err != nil {
		return "", err
	}
	username = strings.TrimSpace(username)

	return username, nil
}

func (manager *chatIOManager) ReadChatRoom() (string, error) {
	fmt.Printf("Enter chat room: ")
	chatRoom, err := manager.inputReader.ReadString('\n')
	if err != nil {
		return "", err
	}
	chatRoom = strings.TrimSpace(chatRoom)

	return chatRoom, nil
}

func (manager *chatIOManager) ReadMessage(username string) (string, error) {
	fmt.Printf("\r[%s] : ", username)
	message, err := manager.inputReader.ReadString('\n')
	if err != nil {
		return "", err
	}
	message = strings.TrimRight(message, "\r\n")

	return message, nil
}

func (manager *chatIOManager) DisplayMessage(receivedMessage *chatv1.Message, username string) {
	fmt.Printf(
		"\r[%s] : %s\n[%s] : ", // replace current line
		receivedMessage.Sender.Username,
		receivedMessage.Content,
		username,
	)
}
