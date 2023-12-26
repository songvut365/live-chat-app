package manager

import (
	"bufio"
	"fmt"
	chatv1 "live-chat-app/app/internal/grpc/gen/chat/v1"
	"os"
	"strings"
)

type IOManager interface {
	ReadUsername() (string, error)
	ReadChatRoom() (string, error)
	ReadMessage(username string) (string, error)
	DisplayMessage(message *chatv1.Message, username string)
}

type ioManager struct {
	reader *bufio.Reader
}

func NewIOManager() IOManager {
	return &ioManager{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (manager *ioManager) ReadUsername() (string, error) {
	fmt.Printf("Enter your username: ")
	username, err := manager.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	username = strings.TrimSpace(username)

	return username, nil
}

func (manager *ioManager) ReadChatRoom() (string, error) {
	fmt.Printf("Enter chat room: ")
	chatRoom, err := manager.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	chatRoom = strings.TrimSpace(chatRoom)

	return chatRoom, nil
}

func (manager *ioManager) ReadMessage(username string) (string, error) {
	fmt.Printf("\r[%s] : ", username)
	message, err := manager.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	message = strings.TrimRight(message, "\r\n")

	return message, nil
}

func (manager *ioManager) DisplayMessage(message *chatv1.Message, username string) {
	fmt.Printf(
		"\r[%s] : %s\n[%s] : ", // replace current line
		message.Sender.Username,
		message.Content,
		username,
	)
}
