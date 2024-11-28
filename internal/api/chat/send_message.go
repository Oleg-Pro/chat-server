package chat

import (
	"context"
//	"log"
	desc "github.com/Oleg-Pro/chat-server/pkg/chat_v1"
	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"		
)

// SendMessage send message
func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*empty.Empty, error) {
	i.mxChannel.RLock()
	chatChannel, ok := i.channels[req.GetChatId()]

	if !ok {
		return nil, status.Errorf(codes.NotFound, "chat not found")
	}


	// Если канал для чата еще не создан
/*	if !ok {
		log.Println("SendMessage creating channel for chat")		
		i.channels[req.GetChatId()] = make(chan *desc.Message, 100)
		chatChannel, ok = i.channels[req.GetChatId()]				
	}


	if !ok {	
		log.Println("Strange things happen!)))")
		return &empty.Empty{}, nil		
	}*/


	i.mxChannel.RUnlock()
	
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
