package chat

import (
	"context"
	"log"

	desc "github.com/Oleg-Pro/chat-server/pkg/chat_v1"
	empty "github.com/golang/protobuf/ptypes/empty"
)

// Delete chat
func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	_, err := i.chatService.Delete(ctx, req.GetId())
	if err != nil {
		log.Printf("Failed to delete user with id %d: %v", req.GetId(), err)
		return nil, err
	}

	return &empty.Empty{}, nil

}
