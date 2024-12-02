package main

import (
	"context"
	"flag"
	"io"
	"log"
	"sync"
	"time"

	desc "github.com/Oleg-Pro/chat-server/pkg/chat_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var accessToken = flag.String("a", "", "access token")

const (
	address = "localhost:50052"
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Failed to close connection %v", err)
		}
	}()

	client := desc.NewChatV1Client(conn)

	ctx := context.Background()
	md := metadata.New(map[string]string{"Authorization": "Bearer " + *accessToken})

	log.Printf("Access token %s\n", *accessToken)
	ctx = metadata.NewOutgoingContext(ctx, md)

	r, err := client.Create(ctx, &desc.CreateRequest{
		UserNames: []string{"user1,", "user55"},
	})
	if err != nil {
		log.Fatalf("Failed to User %+v", err)
	}

	log.Printf(color.RedString("Create response \n"), color.GreenString("%v", r))

	wg := sync.WaitGroup{}
	wg.Add(2)

	chatID := r.GetId()

	go func() {
		defer wg.Done()

		err = connectChat(ctx, client, chatID, "user1", 5*time.Second)
		if err != nil {
			log.Printf("failed to connect chat1: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err = connectChat(ctx, client, chatID, "user2", 7*time.Second)
		if err != nil {
			log.Printf("failed to connect chat2: %v", err)
		}
	}()

	wg.Wait()
	log.Println("Chat client exited")

}

func connectChat(ctx context.Context, client desc.ChatV1Client, chatID int64, username string, period time.Duration) error {
	log.Printf("User %s connecting to chat %d\n", username, chatID)
	stream, err := client.Connect(ctx, &desc.ConnectRequest{
		ChatId:   chatID,
		Username: "oleg",
	})
	if err != nil {
		return err
	}

	go func() {
		for {
			message, errRecv := stream.Recv()
			if errRecv == io.EOF {
				return
			}
			if errRecv != nil {
				log.Println("failed to receive message from stream: ", errRecv)
				return
			}

			log.Printf("[%v] - [from: %s]: %s\n",
				color.YellowString(message.GetCreatedAt().AsTime().Format(time.RFC3339)),
				color.BlueString(message.GetFrom()),
				message.GetText(),
			)
		}
	}()

	for {
		time.Sleep(period)

		text := gofakeit.Word()

		_, err = client.SendMessage(ctx, &desc.SendMessageRequest{
			ChatId: chatID,
			Message: &desc.Message{
				From:      username,
				Text:      text,
				CreatedAt: timestamppb.Now(),
			},
		})
		if err != nil {
			log.Println("failed to send message: ", err)
			return err
		}
	}
}
