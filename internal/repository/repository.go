package repository

import (
	"context"

	"github.com/Oleg-Pro/chat-server/internal/model"
)

// ChatRepository interface for Chat Repository
type ChatRepository interface {
	Create(ctx context.Context, info *model.ChatInfo) (int64, error)
	Delete(ctx context.Context, id int64) (int64, error)
}
