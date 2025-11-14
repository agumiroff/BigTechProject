package app

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/agumiroff/BigTechProject/payment/v1/internal/config"
	"github.com/agumiroff/BigTechProject/payment/v1/migrations"
	"github.com/agumiroff/BigTechProject/platform/closer"
	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/health"
	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
	paymentv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runGRPCServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initLogger,
		a.initCloser,
		a.initDI,
		a.initMongo,
		a.initMigrations,
		a.initListener,
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
	return config.Load()
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDIContainer()
	return nil
}

func (a *App) initMongo(ctx context.Context) error {
	// Инициализируем MongoDB клиент
	_, err := a.diContainer.MongoClient(ctx, config.AppConfig().Mongo.URI())
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Получаем инстанс БД
	db := a.diContainer.MongoDB(config.AppConfig().Mongo.DBName())
	if db == nil {
		return fmt.Errorf("failed to get MongoDB database instance")
	}

	// Добавляем closer для MongoDB
	closer.AddNamed("MongoDB", func(ctx context.Context) error {
		return a.diContainer.Close(ctx)
	})

	logger.Info(ctx, "✅ Connected to MongoDB")

	return nil
}

func (a *App) initMigrations(ctx context.Context) error {
	db := a.diContainer.MongoDB(config.AppConfig().Mongo.DBName())
	if db == nil {
		return fmt.Errorf("MongoDB not initialized")
	}

	migPath := config.AppConfig().Mongo.MigrationPath()
	if err := migrations.ApplyMigrations(ctx, db, migPath); err != nil {
		logger.Warn(ctx, fmt.Sprintf("Warning: Migration error: %v", err))
	}

	logger.Info(ctx, "Migrations applied")

	return nil
}

func (a *App) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().GRPC.Address())
	if err != nil {
		return fmt.Errorf("failed to create TCP listener: %w", err)
	}

	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		lerr := listener.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}
		return nil
	})

	a.listener = listener

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	// Регистрируем reflection для удобства работы с grpcurl
	reflection.Register(a.grpcServer)

	// Регистрируем health service для проверки работоспособности
	health.RegisterService(a.grpcServer)

	// Инициализируем репозиторий
	repo := a.diContainer.PaymentRepository()
	if repo == nil {
		return fmt.Errorf("failed to initialize repository")
	}

	// Инициализируем сервис
	service := a.diContainer.PaymentService()
	if service == nil {
		return fmt.Errorf("failed to initialize service")
	}

	// Получаем gRPC server implementation
	paymentServer := a.diContainer.PaymentServiceServer()
	if paymentServer == nil {
		return fmt.Errorf("failed to initialize payment server")
	}

	// Регистрируем Payment Service
	paymentv1.RegisterPaymentServiceServer(a.grpcServer, paymentServer)

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC PaymentService server listening on %s", config.AppConfig().GRPC.Address()))

	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		return fmt.Errorf("gRPC server error: %w", err)
	}

	return nil
}
