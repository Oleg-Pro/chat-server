package chat

import (
	"context"
	"errors"
	"strings"

	"github.com/Oleg-Pro/chat-server/internal/model"
	desc "github.com/Oleg-Pro/chat-server/pkg/chat_v1"
)

// ErrUserListEmpty error password are not equal
var ErrUserListEmpty = errors.New("passwords are not equal")

// Create implementation of Create User Api Method
func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	if len(req.GetUserNames()) == 0 {
		return nil, ErrUserListEmpty
	}

	users := string(strings.Join(req.GetUserNames(), ","))

	chatID, err := i.chatService.Create(ctx, &model.ChatInfo{Users: users})
	if err != nil {
		return nil, err
	}

	return &desc.CreateResponse{
		Id: chatID,
	}, nil
}
