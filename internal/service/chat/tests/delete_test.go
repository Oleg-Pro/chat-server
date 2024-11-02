package tests

import (
	"context"
	"testing"

	"github.com/Oleg-Pro/chat-server/internal/repository"
	repoMocks "github.com/Oleg-Pro/chat-server/internal/repository/mocks"
	"github.com/Oleg-Pro/chat-server/internal/service/chat"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	t.Parallel()
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository

	type args struct {
		ctx context.Context
		id  int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id           = gofakeit.Int64()
		numberOfRows = int64(1)
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
				id:  id,
			},
			want: numberOfRows,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(numberOfRows, nil)
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
			resonse, err := api.Delete(tt.args.ctx, tt.args.id)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resonse)
		})
	}
}
