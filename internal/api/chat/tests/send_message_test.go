package tests

import (
	"context"
	"database/sql"
	"testing"

	chatAPI "github.com/Oleg-Pro/chat-server/internal/api/chat"
	"github.com/Oleg-Pro/chat-server/internal/model"
	"github.com/Oleg-Pro/chat-server/internal/service"
	serviceMocks "github.com/Oleg-Pro/chat-server/internal/service/mocks"
	desc "github.com/Oleg-Pro/chat-server/pkg/chat_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestSendMessage(t *testing.T) {
	t.Parallel()
	type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		ctx context.Context
		req *desc.SendMessageRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		from        = gofakeit.Name()
		text        = "Text message"
		timestamppb = timestamppb.Now()

		req = &desc.SendMessageRequest{
			From:      from,
			Text:      text,
			Timestamp: timestamppb,
		}

		res = &empty.Empty{}
	)

	defer t.Cleanup(mc.Finish)

	var timestamp sql.NullTime
	if timestamppb == nil {
		timestamp.Valid = false
	} else {
		timestamp.Time = timestamppb.AsTime()
	}

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
				mock.SendMessageMock.Expect(ctx, &model.MessageInfo{
					From:      from,
					Text:      text,
					Timestamp: timestamp,
				}).Return(nil)
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
			resonse, err := api.SendMessage(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resonse)
		})
	}
}
