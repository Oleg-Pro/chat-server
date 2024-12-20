package chat

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Oleg-Pro/chat-server/internal/model"
	"github.com/Oleg-Pro/chat-server/internal/repository"
	"github.com/Oleg-Pro/platform-common/pkg/db"
)

const (
	chatTable = "chats"

	chatColumnID    = "id"
	chatColumnUsers = "users"
)

type repo struct {
	//	pool *pgxpool.Pool
	db db.Client
}

// NewRepository create ChatRepository
func NewRepository(db db.Client) repository.ChatRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, info *model.ChatInfo) (int64, error) {

	builderInsert := sq.Insert(chatTable).
		PlaceholderFormat(sq.Dollar).
		Columns(chatColumnUsers).
		Values(info.Users).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return 0, err
	}

	var userID int64

	q := db.Query{
		Name:     "chat_repository.Create",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&userID)
	if err != nil {
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
		return 0, err
	}

	q := db.Query{
		Name:     "chat_repository.Delete",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected(), nil
}
