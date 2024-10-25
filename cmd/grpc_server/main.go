package main

import (
	"context"
//	"flag"
	"log"
//	"net"

/*	chatAPI "github.com/Oleg-Pro/chat-server/internal/api/chat"
	"github.com/Oleg-Pro/chat-server/internal/config"
	"github.com/Oleg-Pro/chat-server/internal/repository/chat"

	chatService "github.com/Oleg-Pro/chat-server/internal/service/chat"
	desc "github.com/Oleg-Pro/chat-server/pkg/chat_v1"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"*/
	"github.com/Oleg-Pro/chat-server/internal/app"	
)

/*var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}*/

func main() {

	ctx := context.Background()
	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}

/*	flag.Parse()

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

	desc.RegisterChatV1Server(s, chatAPI.NewImplementation(chatService))

	log.Printf("Server listening at %v", listener.Addr())

	if err = s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}*/
}
