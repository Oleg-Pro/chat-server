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
	log.Println("SendMessage RLocking")				
	i.mxChannel.RLock()
	log.Println("SendMessage RLocked")					
	chatChannel, ok := i.channels[req.GetChatId()]	

	// Если канал для чата еще не создан
	if !ok {
		log.Println("SendMessage creating channel for chat")		
		i.channels[req.GetChatId()] = make(chan *desc.Message, 100)
	}

	chatChannel, ok = i.channels[req.GetChatId()]		
	if !ok {	
		log.Println("Strange things happen!)))")
		return &empty.Empty{}, nil		
	}


	log.Println("SendMessage Try to get channel chat")					
	log.Println("SendMessage RUnlocking")						
	i.mxChannel.RUnlock()
	log.Println("SendMessage RUnlocked")					

	//create add chat to channel if it does not exists

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

	log.Println("SendMessage return")
	return &empty.Empty{}, nil
}
