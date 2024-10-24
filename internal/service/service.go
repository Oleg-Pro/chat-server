package service

import (
	"context"

	"github.com/Oleg-Pro/chat-server/internal/model"
)

// ChatService iterface for Chat Service
type ChatService interface {
	Create(ctx context.Context, info *model.ChatInfo) (int64, error)
	Delete(ctx context.Context, id int64) (int64, error)
	SendMessage(ctx context.Context, messageInfo *model.MessageInfo) error
}
