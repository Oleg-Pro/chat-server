package chat

import (
	"context"

	"github.com/Oleg-Pro/chat-server/internal/model"
)

func (s *serv) SendMessage(_ context.Context, _ *model.MessageInfo) error {
	return nil
}
