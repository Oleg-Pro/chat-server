package chat

import (
//	"context"
//	"database/sql"
"log"
//	"github.com/Oleg-Pro/chat-server/internal/model"
	desc "github.com/Oleg-Pro/chat-server/pkg/chat_v1"
//	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
//	"slices"

)

func (i *Implementation) Connect(req *desc.ConnectRequest,stream desc.ChatV1_ConnectServer) error {
	i.mxChannel.RLock()
	chatChan, ok := i.channels[req.GetChatId()]
	i.mxChannel.RUnlock()

	if !ok {
		return status.Errorf(codes.NotFound, "chat not found")
	}

	i.mxChat.Lock()
	if _, okChat := i.chats[req.GetChatId()]; !okChat {
		i.chats[req.GetChatId()] = &Chat{
			streams: make(map[string]desc.ChatV1_ConnectServer),
		}
	}
	i.mxChat.Unlock()	


	log.Printf("Connect: %#v", req)					


	i.chats[req.GetChatId()].m.Lock()
	i.chats[req.GetChatId()].streams[req.GetUsername()] = stream
	i.chats[req.GetChatId()].m.Unlock()


	for {
		select {
		case msg, okCh := <-chatChan:
			if !okCh {
				return nil
			}

			
			for _, stream := range i.chats[req.GetChatId()].streams {
				log.Printf("Connect Sending message to client chatId : %d message : %#v", req.GetChatId(), msg)				
				if err := stream.Send(msg); err != nil {
					return err
				}
			}
		case <-stream.Context().Done():
			// Как понять, какому пользователюя принадленжит стрим и удалить его из мапы?
			// В запросе нет идентификатора пользоватял, непонятно, как удалять из мапы стрим конкретного пользователя

/*			i.chats[req.GetId()].m.Lock()
			// Я пробовал в удалять из слайса так, но не работает
			i.chats[req.GetId()] = slices.DeleteFunc(i.chats[req.GetId()], func(e desc.ChatV1_ConnectServer) bool {
					e == stream
			})
			delete(i.chats[req.GetChatId()].streams, req.GetUsername())
			i.chats[req.GetId()].m.Unlock()*/
			log.Println("Connect Context doen exit")
			return nil
	
		}		
	}



//	stream.Send(&desc.Message{})
	//stream.Send(&desc.Message{C Text: "Test", From: "Oleg"})
//	return nil
	//return status.Errorf(codes.Unimplemented, "method ConnectChat not implemented1")
}

