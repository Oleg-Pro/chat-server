package tests

import (
	"context"
	"testing"

	chatAPI "github.com/Oleg-Pro/chat-server/internal/api/chat"
	"github.com/Oleg-Pro/chat-server/internal/service"
	serviceMocks "github.com/Oleg-Pro/chat-server/internal/service/mocks"
	desc "github.com/Oleg-Pro/chat-server/pkg/chat_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	t.Parallel()
	type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		ctx context.Context
		req *desc.DeleteRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id           = gofakeit.Int64()
		numberOfRows = int64(1)
		req          = &desc.DeleteRequest{
			Id: id,
		}
		res = &empty.Empty{}
	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *empty.Empty
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

				mock.DeleteMock.Expect(ctx, id).Return(numberOfRows, nil)
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
			resonse, err := api.Delete(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resonse)
		})
	}
}
