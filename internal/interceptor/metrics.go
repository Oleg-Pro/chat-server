package interceptor

import (
	"context"
	//	"time"
	"github.com/Oleg-Pro/chat-server/internal/metric"
	"google.golang.org/grpc"
)

// MetricsInterceptor metrics interceptor
func MetricsInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	metric.IncRequestCounter()

	//	timeStart := time.Now()
	res, err := handler(ctx, req)
	/*if err != nil {
		logger.Error(error.Error())
	}*/

	return res, err

}
