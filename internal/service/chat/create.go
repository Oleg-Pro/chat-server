package chat

import (
	"context"

	"github.com/Oleg-Pro/chat-server/internal/model"
)

func (s *serv) Create(ctx context.Context, info *model.ChatInfo) (int64, error) {
	return s.chatRepository.Create(ctx, info)
}
