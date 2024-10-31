package tests

import (
	"context"
	"testing"

	chatAPI "github.com/Oleg-Pro/chat-server/internal/api/chat"
	"github.com/Oleg-Pro/chat-server/internal/model"
	"github.com/Oleg-Pro/chat-server/internal/service"
	serviceMocks "github.com/Oleg-Pro/chat-server/internal/service/mocks"
	desc "github.com/Oleg-Pro/chat-server/pkg/chat_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		userNames = []string{"user1", "user2", "user3"}

		req = &desc.CreateRequest{
			UserNames: userNames,
		}

		emptyUserNames        = []string{}
		reqWithEmptyUserNames = &desc.CreateRequest{
			UserNames: emptyUserNames,
		}

		chatInfo = &model.ChatInfo{
			Users: "user1,user2,user3",
		}

		res = &desc.CreateResponse{
			Id: id,
		}
	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateResponse
		err             error
		chatServiceMock chatServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)

				mock.CreateMock.Expect(ctx, chatInfo).Return(id, nil)
				return mock
			},
		},
		{
			name: "empty user list",
			args: args{
				ctx: ctx,
				req: reqWithEmptyUserNames,
			},
			want: nil,
			err:  chatAPI.ErrUserListEmpty,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			chatServiceMock := tt.chatServiceMock(mc)
			api := chatAPI.NewImplementation(chatServiceMock)
			resonse, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resonse)
		})
	}
}
