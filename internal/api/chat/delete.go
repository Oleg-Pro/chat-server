package chat

import (
	"context"

	desc "github.com/Oleg-Pro/chat-server/pkg/chat_v1"
	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/opentracing/opentracing-go"
)

// Delete chat
func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "delete chat")
	defer span.Finish()
	_, err := i.chatService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil

}
