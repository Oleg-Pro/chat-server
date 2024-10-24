package chat

import (
	"context"
	"database/sql"
	"log"

	"github.com/Oleg-Pro/chat-server/internal/model"
	desc "github.com/Oleg-Pro/chat-server/pkg/chat_v1"
	empty "github.com/golang/protobuf/ptypes/empty"
)

// SendMessage send message
func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*empty.Empty, error) {
	log.Printf("Send req: %+v", req)

	var timestamp sql.NullTime
	if req.GetTimestamp() == nil {
		timestamp.Valid = false

	} else {
		timestamp.Time = req.GetTimestamp().AsTime()
	}

	err := i.chatService.SendMessage(ctx, &model.MessageInfo{
		From:      req.GetFrom(),
		Text:      req.GetText(),
		Timestamp: timestamp,
	})

	if err != nil {
		log.Printf("Failed to send message: %s", err.Error())
		return nil, err
	}

	return &empty.Empty{}, nil
}
