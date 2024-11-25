package app

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Oleg-Pro/chat-server/internal/config"
	"github.com/Oleg-Pro/chat-server/internal/interceptor"
	"github.com/Oleg-Pro/chat-server/internal/logger"
	"github.com/Oleg-Pro/chat-server/internal/metric"
	desc "github.com/Oleg-Pro/chat-server/pkg/chat_v1"
	"github.com/Oleg-Pro/platform-common/pkg/closer"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/natefinch/lumberjack"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// App type
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	configPath      string
	logLevel        string
}

// NewApp creats App
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}
	flag.StringVar(&a.configPath, "config-path", ".env", "path to config file")
	a.logLevel = *flag.String("l", "info", "log level")
	flag.Parse()

	logger.Init(a.getCore(a.getAtomicLevel()))
	err := metric.Init(ctx)
	if err != nil {
		log.Fatalf("failed to init metrics: %v", err)
	}

	err = a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run runs App
func (a *App) Run() error {

	defer func() {
		closer.CloseAll()
		closer.Wait()

	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := a.runGRPCServer()
		if err != nil {
			log.Fatalf("Failed to run grpc server : %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		err := a.runPrometheusServer()
		if err != nil {
			log.Fatalf("Failed to run prometheus server : %v", err)
		}
	}()

	wg.Wait()
	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(a.configPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				interceptor.MetricsInterceptor,
				interceptor.LogInterceptor,
				a.serviceProvider.AuthInterceptor(ctx).AcccessInterceptor,
			),
		),
	)

	reflection.Register(a.grpcServer)
	desc.RegisterChatV1Server(a.grpcServer, a.serviceProvider.ChatImplementation(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	logger.Info(fmt.Sprintf("GRPC server is running on %s", a.serviceProvider.GRPCConfig().Address()))
	listener, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())

	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	return nil
}

func (a *App) runPrometheusServer() error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	prometheusServer := &http.Server{
		Addr:        a.serviceProvider.PrometheusServerConfig().Address(),
		Handler:     mux,
		ReadTimeout: 3 * time.Second,
	}

	logger.Info(fmt.Sprintf("Prometheus server is running on %s", a.serviceProvider.PrometheusServerConfig().Address()))

	err := prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) getCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)
}

func (a *App) getAtomicLevel() zap.AtomicLevel {
	var level zapcore.Level
	if err := level.Set(a.logLevel); err != nil {
		log.Fatalf("failed to set log level: %v", err)
	}

	return zap.NewAtomicLevelAt(level)
}
