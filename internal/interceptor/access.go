package interceptor

import (
	"context"
	"fmt"
	"log"
	"strings"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"		
	"github.com/opentracing/opentracing-go"		
	accessDesc "github.com/Oleg-Pro/auth/pkg/access_v1"
	"github.com/Oleg-Pro/chat-server/internal/logger"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	authPrefix = "Bearer "
	authPort   = 50051
)

// AuthInterceptor auth interceptor struct
type AuthInterceptor struct {
	AccessV1Client accessDesc.AccessV1Client
}

// AcccessInterceptor access interceptor
func (authInterceptor AuthInterceptor) AcccessInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	logger.Info(fmt.Sprintf("Incerceptor FullMethod : %s\n", info.FullMethod))

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not provided")
	}

	logger.Info(fmt.Sprintf("MD : %#v\n", md))

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return nil, errors.New("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)
	logger.Info(fmt.Sprintf("AccesToken : %#v\n", accessToken))
	clientCtx := context.Background()
	clientCtx = metadata.NewOutgoingContext(clientCtx, md)

	span, clientCtx := opentracing.StartSpanFromContext(clientCtx, "authorization")

	defer span.Finish()

	_, err := authInterceptor.AccessV1Client.Check(clientCtx, &accessDesc.CheckRequest{
		EndpointAddress: info.FullMethod,
	})
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

// NewAuthInterceptor AuthInterceptor constructor
func NewAuthInterceptor() *AuthInterceptor {

	conn, err := grpc.Dial(
		fmt.Sprintf(":%d", authPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),		
	)

	if err != nil {
		log.Fatalf("Error create grpc connect to auth: %s", err.Error())
	}

	cl := accessDesc.NewAccessV1Client(conn)

	return &AuthInterceptor{AccessV1Client: cl}
}
