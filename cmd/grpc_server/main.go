package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"database/sql"
	"github.com/Oleg-Pro/chat-server/internal/config"
	"github.com/Oleg-Pro/chat-server/internal/model"
	"github.com/Oleg-Pro/chat-server/internal/repository/chat"
	"github.com/Oleg-Pro/chat-server/internal/service"
	chatService "github.com/Oleg-Pro/chat-server/internal/service/chat"
	desc "github.com/Oleg-Pro/chat-server/pkg/chat_v1"
	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedChatV1Server
	chatService    service.ChatService
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	if len(req.GetUserNames()) == 0 {
		err := fmt.Errorf("users list should not be empty")
		log.Printf("Create Chat Error: %v", err)

		return nil, err
	}

	users := string(strings.Join(req.GetUserNames(), ","))

	chatID, err := s.chatService.Create(ctx, &model.ChatInfo{Users: users})
	if err != nil {
		log.Printf("Failed to create chat: %v", err)
		return nil, err
	}

	return &desc.CreateResponse{
		Id: chatID,
	}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	_, err := s.chatService.Delete(ctx, req.GetId())
	if err != nil {
		log.Printf("Failed to delete user with id %d: %v", req.GetId(), err)
		return nil, err
	}

	return &empty.Empty{}, nil

}

func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*empty.Empty, error) {
	log.Printf("Send req: %+v", req)

	var timestamp sql.NullTime
	if (req.GetTimestamp() == nil) {
		timestamp.Valid = false

	} else {
		timestamp.Time = req.GetTimestamp().AsTime()
	}


	err := s.chatService.SendMessage(ctx, &model.MessageInfo{
		From: req.GetFrom(),
		Text: req.GetText(),
		Timestamp: timestamp,

	})

	if err != nil {
		log.Printf("Failed to send message: %s", err.Error())
		return nil, err
	}

	return &empty.Empty{}, nil
}

func main() {
	flag.Parse()

	// Считываем переменные окружения
	log.Printf("confiPath :%s", configPath)
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Println("Loaded Config Parameters")

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	listener, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	chatRepository := chat.NewRepository(pool)
	chatService := chatService.New(chatRepository)

	desc.RegisterChatV1Server(s, &server{chatService: chatService})
	log.Printf("Server listening at %v", listener.Addr())

	if err = s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
