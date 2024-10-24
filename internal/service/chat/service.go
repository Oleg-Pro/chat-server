package chat

import (
	"github.com/Oleg-Pro/chat-server/internal/repository"
	"github.com/Oleg-Pro/chat-server/internal/service"
)

type serv struct {
	chatRepository repository.ChatRepository
}

// New create ChatService
func New(chatRepository repository.ChatRepository) service.ChatService {
	return &serv{chatRepository: chatRepository}
}
