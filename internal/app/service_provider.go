package app

import (
	"context"
	"log"

	chatAPI "github.com/Oleg-Pro/chat-server/internal/api/chat"
	"github.com/Oleg-Pro/chat-server/internal/config"
	"github.com/Oleg-Pro/chat-server/internal/interceptor"
	"github.com/Oleg-Pro/chat-server/internal/repository"
	chatRepository "github.com/Oleg-Pro/chat-server/internal/repository/chat"
	"github.com/Oleg-Pro/chat-server/internal/service"
	chatService "github.com/Oleg-Pro/chat-server/internal/service/chat"
	"github.com/Oleg-Pro/platform-common/pkg/closer"
	"github.com/Oleg-Pro/platform-common/pkg/db"
	"github.com/Oleg-Pro/platform-common/pkg/db/pg"
	"github.com/Oleg-Pro/platform-common/pkg/db/transaction"
)

type serviceProvider struct {
	pgConfig               config.PGConfig
	prometheusServerConfig config.PrometheusServerConfig
	grpcConfig             config.GRPCConfig

	dbClient       db.Client
	txManager      db.TxManager
	chatRepository repository.ChatRepository

	chatService       service.ChatService
	chatImplemenation *chatAPI.Implementation

	authInterceptor *interceptor.AuthInterceptor
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) PrometheusServerConfig() config.PrometheusServerConfig {
	if s.prometheusServerConfig == nil {
		cfg, err := config.NewPrometheusServerConfig()
		if err != nil {
			log.Fatalf("failed to get prometheus server config: %s", err.Error())
		}

		s.prometheusServerConfig = cfg
	}

	return s.prometheusServerConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		client, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = client.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(client.Close)

		s.dbClient = client
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}
	return s.txManager
}

func (s *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = chatRepository.NewRepository(s.DBClient(ctx))
	}

	return s.chatRepository
}

func (s *serviceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.New(s.ChatRepository(ctx))
	}

	return s.chatService
}

func (s *serviceProvider) ChatImplementation(ctx context.Context) *chatAPI.Implementation {

	if s.chatImplemenation == nil {
		s.chatImplemenation = chatAPI.NewImplementation(s.ChatService(ctx))
	}
	return s.chatImplemenation
}

func (s *serviceProvider) AuthInterceptor(_ context.Context) *interceptor.AuthInterceptor {
	if s.authInterceptor == nil {
		s.authInterceptor = interceptor.NewAuthInterceptor()
	}

	return s.authInterceptor
}
