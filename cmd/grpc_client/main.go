package main

import (
	"context"
	"log"
	"time"

	desc "github.com/Oleg-Pro/chat-server/pkg/chat_v1"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50502"
)

func main() {
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	r, err := client.Create(ctx, &desc.CreateRequest{
		UserNames: []string{"user1,", "user2"},
	})
	if err != nil {
		log.Fatalf("Failed to User %+v", err)
	}

	log.Printf(color.RedString("Create response \n"), color.GreenString("%v", r))

}
