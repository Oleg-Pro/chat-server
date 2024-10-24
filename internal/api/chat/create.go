package chat

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Oleg-Pro/chat-server/internal/model"
	desc "github.com/Oleg-Pro/chat-server/pkg/chat_v1"
)

// Create implementation of Create User Api Method
func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	if len(req.GetUserNames()) == 0 {
		err := fmt.Errorf("users list should not be empty")
		log.Printf("Create Chat Error: %v", err)

		return nil, err
	}

	users := string(strings.Join(req.GetUserNames(), ","))

	chatID, err := i.chatService.Create(ctx, &model.ChatInfo{Users: users})
	if err != nil {
		log.Printf("Failed to create chat: %v", err)
		return nil, err
	}

	return &desc.CreateResponse{
		Id: chatID,
	}, nil
}
