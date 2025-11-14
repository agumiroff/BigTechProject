package app

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	api "github.com/agumiroff/BigTechProject/inventory/v1/internal/api/inventory/v1"
	"github.com/agumiroff/BigTechProject/inventory/v1/internal/repository"
	partRepo "github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/part"
	"github.com/agumiroff/BigTechProject/inventory/v1/internal/service"
	"github.com/agumiroff/BigTechProject/inventory/v1/internal/service/part"
	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
	inventoryv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	client       *mongo.Client
	db           *mongo.Database
	repo         repository.InvRepository
	service      service.InvService
	invAPIServer inventoryv1.InventoryServiceServer
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

// MongoClient инициализирует и возвращает MongoDB клиент
func (d *diContainer) MongoClient(ctx context.Context, uri string) (*mongo.Client, error) {
	if d.client == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			logger.Error(ctx, "failed to connect mongodb",
				zap.String("uri", uri),
				zap.Error(err),
			)
			return nil, err
		}

		// Проверяем подключение
		if err = client.Ping(ctx, nil); err != nil {
			logger.Error(ctx, "mongodb ping failed",
				zap.Error(err),
			)
			return nil, err
		}

		d.client = client
	}

	return d.client, nil
}

// MongoDB возвращает инстанс MongoDB базы данных
func (d *diContainer) MongoDB(dbName string) *mongo.Database {
	if d.db == nil && d.client != nil {
		d.db = d.client.Database(dbName)
	}
	return d.db
}

// InvRepository инициализирует и возвращает репозиторий
func (d *diContainer) InvRepository(ctx context.Context) repository.InvRepository {
	if d.repo == nil && d.db != nil {
		d.repo = partRepo.NewRepository(ctx, d.db)
	}
	return d.repo
}

// InvService инициализирует и возвращает сервис
func (d *diContainer) InvService() service.InvService {
	if d.service == nil && d.repo != nil {
		d.service = part.NewService(d.repo)
	}
	return d.service
}

// InventoryServiceServer инициализирует и возвращает gRPC server implementation
func (d *diContainer) InventoryServiceServer() inventoryv1.InventoryServiceServer {
	if d.invAPIServer == nil && d.service != nil {
		d.invAPIServer = api.NewAPI(d.service)
	}
	return d.invAPIServer
}

// Close закрывает все соединения
func (d *diContainer) Close(ctx context.Context) error {
	if d.client != nil {
		return d.client.Disconnect(ctx)
	}
	return nil
}
