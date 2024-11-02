package tests

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Oleg-Pro/chat-server/internal/model"
	"github.com/Oleg-Pro/chat-server/internal/repository"
	repoMocks "github.com/Oleg-Pro/chat-server/internal/repository/mocks"
	"github.com/Oleg-Pro/chat-server/internal/service/chat"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestMessage(t *testing.T) {
	t.Parallel()
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository

	type args struct {
		ctx context.Context
		req *model.MessageInfo
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		from = gofakeit.Name()
		text = "Text message"

		timestamp = sql.NullTime{Time: time.Now(), Valid: true}

		req = &model.MessageInfo{
			From:      from,
			Text:      text,
			Timestamp: timestamp,
		}
	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		userRepositoryMock userRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepoMock := tt.userRepositoryMock(mc)
			api := chat.New(userRepoMock)
			err := api.SendMessage(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}
}
