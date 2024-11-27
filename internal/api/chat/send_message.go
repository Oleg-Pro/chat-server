package chat

import (
	"context"

//	"database/sql"


	//"github.com/Oleg-Pro/chat-server/internal/logger"
	"log"
	desc "github.com/Oleg-Pro/chat-server/pkg/chat_v1"
	empty "github.com/golang/protobuf/ptypes/empty"
)

// SendMessage send message
func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*empty.Empty, error) {
	i.mxChannel.RLock()
	chatChannel, ok := i.channels[req.GetChatId()]
	i.mxChannel.RUnlock()

	//create add chat to channel if it does not exists
	if !ok {
		i.channels[req.GetChatId()] = make(chan *desc.Message, 100)
	}

	log.Printf("Send chatId %d  message: %#v\n", req.ChatId, req.GetMessage())

	

	chatChannel <- req.GetMessage()


/*	var timestamp sql.NullTime
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
		return nil, err
	}*/

	return &empty.Empty{}, nil
}
