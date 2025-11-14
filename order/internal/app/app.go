package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/agumiroff/BigTechProject/order/v1/internal/config"
	"github.com/agumiroff/BigTechProject/platform/closer"
	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
	orderV1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
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
	return a.runHTTPServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initLogger,
		a.initCloser,
		a.initDI,
		a.initPostgres,
		a.initMigrations,
		a.initGRPCClients,
		a.initHTTPServer,
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

func (a *App) initPostgres(ctx context.Context) error {
	// Инициализируем PostgreSQL клиент
	_, err := a.diContainer.PostgresDB(ctx, config.AppConfig().Postgres.DSN())
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Добавляем closer для PostgreSQL
	closer.AddNamed("PostgreSQL", func(ctx context.Context) error {
		return a.diContainer.Close(ctx)
	})

	logger.Info(ctx, "✅ Connected to PostgreSQL")

	return nil
}

func (a *App) initMigrations(ctx context.Context) error {
	migPath := config.AppConfig().Postgres.MigPath()
	if err := a.diContainer.RunMigrations(migPath); err != nil {
		logger.Warn(ctx, fmt.Sprintf("Warning: Migration error: %v", err))
	}

	logger.Info(ctx, "Migrations applied")

	return nil
}

func (a *App) initGRPCClients(ctx context.Context) error {
	// Подключаемся к Inventory сервису
	_, err := a.diContainer.InventoryGRPCClient(ctx, config.AppConfig().HTTP.InventoryAddr())
	if err != nil {
		return fmt.Errorf("failed to connect to Inventory service: %w", err)
	}

	// Подключаемся к Payment сервису
	_, err = a.diContainer.PaymentGRPCClient(ctx, config.AppConfig().HTTP.PaymentAddr())
	if err != nil {
		return fmt.Errorf("failed to connect to Payment service: %w", err)
	}

	logger.Info(ctx, "✅ Connected to gRPC services (Inventory, Payment)")

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	// Инициализируем все слои через DI
	_ = a.diContainer.InternalRepository()
	_ = a.diContainer.ExternalRepository()
	_ = a.diContainer.OrderService()
	_ = a.diContainer.OrderAPI()

	// Получаем HTTP handler
	handler := a.diContainer.OrderHandler()
	if handler == nil {
		return fmt.Errorf("failed to initialize order handler")
	}

	// Создаём OpenAPI сервер
	orderServer, err := orderV1.NewServer(handler)
	if err != nil {
		return fmt.Errorf("failed to create OpenAPI server: %w", err)
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчики OpenAPI
	r.Mount("/", orderServer)

	// Создаём HTTP сервер
	a.httpServer = &http.Server{
		Addr:              config.AppConfig().HTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: config.AppConfig().HTTP.Timeout(),
	}

	// Добавляем closer для HTTP сервера
	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		return a.httpServer.Shutdown(shutdownCtx)
	})

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 HTTP OrderService server listening on %s", config.AppConfig().HTTP.Address()))

	err := a.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server error: %w", err)
	}

	return nil
}
