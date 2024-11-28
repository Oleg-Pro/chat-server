package chat

import (
	"sync"	
	"github.com/Oleg-Pro/chat-server/internal/service"
	desc "github.com/Oleg-Pro/chat-server/pkg/chat_v1"
)

type Chat struct {
	streams []desc.ChatV1_ConnectServer
	m       sync.RWMutex
}


// Implementation implementation of Chat API
type Implementation struct {
	desc.UnimplementedChatV1Server
	chatService service.ChatService

	chats  map[int64]*Chat
	mxChat sync.RWMutex

	channels  map[int64]chan *desc.Message
	mxChannel sync.RWMutex
}


func (i *Implementation) CreateChatChannel(id int64) {
	i.channels[id] = make(chan *desc.Message, 100)	
}

// NewImplementation create Chat Api implementation
func NewImplementation(chatService service.ChatService) *Implementation {
	return &Implementation{
		chatService: chatService,
		chats:    make(map[int64]*Chat),
		channels: make(map[int64]chan *desc.Message),
	}
}
