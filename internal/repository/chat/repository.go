package chat

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/Oleg-Pro/chat-server/internal/model"
	"github.com/Oleg-Pro/chat-server/internal/repository"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	chatTable = "chats"

	chatColumnID    = "id"
	chatColumnUsers = "users"
)

type repo struct {
	pool *pgxpool.Pool
}

// NewRepository create UserRepository
func NewRepository(pool *pgxpool.Pool) repository.ChatRepository {
	return &repo{pool: pool}
}

func (r *repo) Create(ctx context.Context, info *model.ChatInfo) (int64, error) {

	builderInsert := sq.Insert(chatTable).
		PlaceholderFormat(sq.Dollar).
		Columns(chatColumnUsers).
		Values(info.Users).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("Failed to build insert query: %v", err)
		return 0, err
	}

	var userID int64

	err = r.pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return 0, err
	}

	return userID, nil
}

func (r *repo) Delete(ctx context.Context, id int64) (int64, error) {

	builderDelete := sq.Delete(chatTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{fmt.Sprintf(`"%s"`, chatColumnID): id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Printf("Failed to build delete query: %v", err)
		return 0, err
	}

	res, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("Failed to delete user with id %d: %v", id, err)
		return 0, err
	}

	return res.RowsAffected(), nil
}
