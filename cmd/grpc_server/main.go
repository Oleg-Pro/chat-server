package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	sq "github.com/Masterminds/squirrel"
	desc "github.com/Oleg-Pro/chat-server/pkg/chat_v1"
	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50502

const (
	dbDSN      = "host=localhost port=54322 dbname=chat-server user=chat-server-user password=chat-server-password sslmode=disable"
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

	fmt.Printf("SQL Query: %s\n", query)

	fmt.Printf("Args: %v %T\n", args, args)
	fmt.Printf("Users: %v %T\n", users, users)

	var chatID int64

	//Почему паника возникает здесь?
	err = s.pool.QueryRow(ctx, query, args...).Scan(&chatID)

	if err != nil {
		log.Printf("Failed to create chat: %v", err)
		return &desc.CreateResponse{}, err
	}

	return &desc.CreateResponse{
		Id: chatID,
	}, nil
}

func (s *server) Delete(_ context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	log.Printf("Delete req: %+v", req)

	return &empty.Empty{}, nil
}

func (s *server) SendMessage(_ context.Context, req *desc.SendMessageRequest) (*empty.Empty, error) {
	log.Printf("Send req: %+v", req)

	return &empty.Empty{}, nil
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)

	}

	s := grpc.NewServer()
	reflection.Register(s)

	desc.RegisterChatV1Server(s, &server{})
	log.Printf("Server listening at %v", listener.Addr())

	if err = s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
