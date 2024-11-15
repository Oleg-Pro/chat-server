package interceptor

import (
	"context"
	"fmt"
	"log"
	"strings"

	accessDesc "github.com/Oleg-Pro/auth/pkg/access_v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	authPrefix = "Bearer "
	authPort   = 50051
)

// AcccessInterceptor access interceptor
func AcccessInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Incerceptor FullMethod : %s\n", info.FullMethod)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not provided")
	}

	log.Printf("MD : %#v\n", md)

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return nil, errors.New("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)
	log.Printf("AccesToken : %#v\n", accessToken)

	clientCtx := context.Background()
	//	md := metadata.New(map[string]string{"Authorization": "Bearer " + accessToken})
	clientCtx = metadata.NewOutgoingContext(clientCtx, md)

	conn, err := grpc.Dial(
		fmt.Sprintf(":%d", authPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	cl := accessDesc.NewAccessV1Client(conn)

	_, err = cl.Check(clientCtx, &accessDesc.CheckRequest{
		EndpointAddress: info.FullMethod,
	})
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}
