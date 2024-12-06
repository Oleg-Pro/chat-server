package chat

import (
	"context"

	desc "github.com/Oleg-Pro/chat-server/pkg/chat_v1"
	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SendMessage send message
func (i *Implementation) SendMessage(_ context.Context, req *desc.SendMessageRequest) (*empty.Empty, error) {
	i.mxChannel.RLock()
	chatChannel, ok := i.channels[req.GetChatId()]

	if !ok {
		return nil, status.Errorf(codes.NotFound, "chat not found")
	}

	i.mxChannel.RUnlock()

	chatChannel <- req.GetMessage()

	return &empty.Empty{}, nil
}
