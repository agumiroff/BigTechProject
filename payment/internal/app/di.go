package app

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	api "github.com/agumiroff/BigTechProject/payment/v1/internal/api/v1"
	"github.com/agumiroff/BigTechProject/payment/v1/internal/repository"
	paymentRepo "github.com/agumiroff/BigTechProject/payment/v1/internal/repository/payment"
	"github.com/agumiroff/BigTechProject/payment/v1/internal/service"
	"github.com/agumiroff/BigTechProject/payment/v1/internal/service/payment"
	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
	paymentv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	client           *mongo.Client
	db               *mongo.Database
	repo             repository.PaymentRepository
	service          service.PaymentService
	paymentAPIServer paymentv1.PaymentServiceServer
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

// PaymentRepository инициализирует и возвращает репозиторий
func (d *diContainer) PaymentRepository() repository.PaymentRepository {
	if d.repo == nil && d.db != nil {
		d.repo = paymentRepo.NewRepository(d.db)
	}
	return d.repo
}

// PaymentService инициализирует и возвращает сервис
func (d *diContainer) PaymentService() service.PaymentService {
	if d.service == nil && d.repo != nil {
		d.service = payment.NewService(d.repo)
	}
	return d.service
}

// PaymentServiceServer инициализирует и возвращает gRPC server implementation
func (d *diContainer) PaymentServiceServer() paymentv1.PaymentServiceServer {
	if d.paymentAPIServer == nil && d.service != nil {
		d.paymentAPIServer = api.NewAPI(d.service)
	}
	return d.paymentAPIServer
}

// Close закрывает все соединения
func (d *diContainer) Close(ctx context.Context) error {
	if d.client != nil {
		return d.client.Disconnect(ctx)
	}
	return nil
}
