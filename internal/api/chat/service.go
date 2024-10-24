package chat

import (
	"github.com/Oleg-Pro/chat-server/internal/service"
	desc "github.com/Oleg-Pro/chat-server/pkg/chat_v1"
)

// Implementation implementation of Chat API
type Implementation struct {
	desc.UnimplementedChatV1Server
	chatService service.ChatService
}

// NewImplementation create Chat Api implementation
func NewImplementation(chatService service.ChatService) *Implementation {
	return &Implementation{chatService: chatService}
}
