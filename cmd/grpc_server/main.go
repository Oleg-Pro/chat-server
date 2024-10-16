package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/Oleg-Pro/chat-server/internal/config"
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

const (
	chatsTable = "chats"
)

type server struct {
	desc.UnimplementedChatV1Server
	pool *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("Create req: %+v", req)

	if len(req.GetUserNames()) == 0 {
		err := fmt.Errorf("users list should not be empty")
		log.Printf("Create Chat Error: %v", err)

		return &desc.CreateResponse{}, err
	}

	users := string(strings.Join(req.GetUserNames(), ","))

	builderInsert := sq.Insert(chatsTable).
		PlaceholderFormat(sq.Dollar).
		Columns("users").
		Values(users).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("Failed to build insert query: %v", err)

		return &desc.CreateResponse{}, err
	}


	var chatID int64

	err = s.pool.QueryRow(ctx, query, args...).Scan(&chatID)

	if err != nil {
		log.Printf("Failed to create chat: %v", err)
		return &desc.CreateResponse{}, err
	}

	return &desc.CreateResponse{
		Id: chatID,
	}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	log.Printf("Delete req: %+v", req)

	builderDelete := sq.Delete(chatsTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Printf("Failed to build delete query: %v", err)
		return &empty.Empty{}, err
	}


	res, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("Failed to delete user with id %d: %v", req.GetId(), err)
		return &empty.Empty{}, err
	}

	log.Printf("delete %d rows", res.RowsAffected())

	return &empty.Empty{}, nil
}

func (s *server) SendMessage(_ context.Context, req *desc.SendMessageRequest) (*empty.Empty, error) {
	log.Printf("Send req: %+v", req)

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

	desc.RegisterChatV1Server(s, &server{pool: pool})
	log.Printf("Server listening at %v", listener.Addr())

	if err = s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
