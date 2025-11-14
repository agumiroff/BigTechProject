package app

import (
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	extRepo "github.com/agumiroff/BigTechProject/order/v1/external/repository"
	extRepoOrder "github.com/agumiroff/BigTechProject/order/v1/external/repository/order"
	"github.com/agumiroff/BigTechProject/order/v1/internal/api"
	apiV1 "github.com/agumiroff/BigTechProject/order/v1/internal/api/v1"
	"github.com/agumiroff/BigTechProject/order/v1/internal/db"
	"github.com/agumiroff/BigTechProject/order/v1/internal/handler/order"
	"github.com/agumiroff/BigTechProject/order/v1/internal/migrator"
	intRepo "github.com/agumiroff/BigTechProject/order/v1/internal/repository"
	intRepoOrder "github.com/agumiroff/BigTechProject/order/v1/internal/repository/order"
	"github.com/agumiroff/BigTechProject/order/v1/internal/service"
	orderService "github.com/agumiroff/BigTechProject/order/v1/internal/service/order"
	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
)

type diContainer struct {
	db            *sql.DB
	inventoryConn *grpc.ClientConn
	paymentConn   *grpc.ClientConn
	internalRepo  intRepo.OrderRepository
	externalRepo  extRepo.OrderRepository
	service       service.OrderService
	api           api.OrderAPI
	handler       *order.OrderHandler
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

// PostgresDB инициализирует и возвращает PostgreSQL подключение
func (d *diContainer) PostgresDB(ctx context.Context, dsn string) (*sql.DB, error) {
	if d.db == nil {
		database, err := db.ConnectDB(dsn)
		if err != nil {
			logger.Error(ctx, "failed to connect to PostgreSQL",
				zap.Error(err),
			)
			return nil, err
		}

		// Проверяем подключение
		if err = database.Ping(); err != nil {
			logger.Error(ctx, "PostgreSQL ping failed",
				zap.Error(err),
			)
			return nil, err
		}

		d.db = database
	}

	return d.db, nil
}

// RunMigrations применяет миграции
func (d *diContainer) RunMigrations(migrationsDir string) error {
	if d.db == nil {
		return fmt.Errorf("database not initialized")
	}

	m := migrator.NewMigrator(d.db, migrationsDir)
	if err := m.RunMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// InventoryGRPCClient инициализирует подключение к Inventory сервису
func (d *diContainer) InventoryGRPCClient(ctx context.Context, address string) (*grpc.ClientConn, error) {
	if d.inventoryConn == nil {
		conn, err := grpc.NewClient(
			address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			logger.Error(ctx, "failed to connect to Inventory service",
				zap.String("address", address),
				zap.Error(err),
			)
			return nil, err
		}

		d.inventoryConn = conn
	}

	return d.inventoryConn, nil
}

// PaymentGRPCClient инициализирует подключение к Payment сервису
func (d *diContainer) PaymentGRPCClient(ctx context.Context, address string) (*grpc.ClientConn, error) {
	if d.paymentConn == nil {
		conn, err := grpc.NewClient(
			address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			logger.Error(ctx, "failed to connect to Payment service",
				zap.String("address", address),
				zap.Error(err),
			)
			return nil, err
		}

		d.paymentConn = conn
	}

	return d.paymentConn, nil
}

// InternalRepository инициализирует и возвращает внутренний репозиторий
func (d *diContainer) InternalRepository() intRepo.OrderRepository {
	if d.internalRepo == nil && d.db != nil {
		d.internalRepo = intRepoOrder.NewRepository(d.db)
	}
	return d.internalRepo
}

// ExternalRepository инициализирует и возвращает внешний репозиторий (gRPC клиенты)
func (d *diContainer) ExternalRepository() extRepo.OrderRepository {
	if d.externalRepo == nil && d.inventoryConn != nil && d.paymentConn != nil {
		d.externalRepo = extRepoOrder.NewRepository(d.inventoryConn, d.paymentConn)
	}
	return d.externalRepo
}

// OrderService инициализирует и возвращает сервис
func (d *diContainer) OrderService() service.OrderService {
	if d.service == nil && d.internalRepo != nil && d.externalRepo != nil {
		d.service = orderService.NewService(d.internalRepo, d.externalRepo)
	}
	return d.service
}

// OrderAPI инициализирует и возвращает API
func (d *diContainer) OrderAPI() api.OrderAPI {
	if d.api == nil && d.service != nil {
		d.api = apiV1.NewAPI(d.service)
	}
	return d.api
}

// OrderHandler инициализирует и возвращает HTTP handler
func (d *diContainer) OrderHandler() *order.OrderHandler {
	if d.handler == nil && d.api != nil {
		d.handler = order.NewHandler(d.api)
	}
	return d.handler
}

// Close закрывает все соединения
func (d *diContainer) Close(ctx context.Context) error {
	var errs []error

	if d.inventoryConn != nil {
		if err := d.inventoryConn.Close(); err != nil {
			logger.Error(ctx, "failed to close Inventory gRPC connection", zap.Error(err))
			errs = append(errs, err)
		}
	}

	if d.paymentConn != nil {
		if err := d.paymentConn.Close(); err != nil {
			logger.Error(ctx, "failed to close Payment gRPC connection", zap.Error(err))
			errs = append(errs, err)
		}
	}

	if d.db != nil {
		if err := d.db.Close(); err != nil {
			logger.Error(ctx, "failed to close PostgreSQL connection", zap.Error(err))
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors during close: %v", errs)
	}

	return nil
}
