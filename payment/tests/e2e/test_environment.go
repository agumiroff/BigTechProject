package integration

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"go.mongodb.org/mongo-driver/bson"

	paymentv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

// InsertTestPayment — вставляет тестовый платеж в коллекцию Mongo и возвращает его transaction UUID
func (env *TestEnvironment) InsertTestPayment(ctx context.Context) (string, error) {
	transactionUUID := gofakeit.UUID()
	now := time.Now().Unix()

	paymentDoc := bson.M{
		"transaction_uuid": transactionUUID,
		"order_uuid":       gofakeit.UUID(),
		"user_uuid":        gofakeit.UUID(),
		"payment_method":   int32(paymentv1.PaymentMethod_PAYMENT_METHOD_CARD),
		"status":           "completed",
		"created_at":       now,
		"updated_at":       now,
	}

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_INITDB_DATABASE")
	if databaseName == "" {
		databaseName = "payments" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(paymentsCollectionName).InsertOne(ctx, paymentDoc)
	if err != nil {
		return "", err
	}

	return transactionUUID, nil
}

// InsertTestPaymentWithData — вставляет тестовый платеж с заданными данными
func (env *TestEnvironment) InsertTestPaymentWithData(ctx context.Context, payment *paymentv1.Payment) (string, error) {
	transactionUUID := gofakeit.UUID()
	now := time.Now().Unix()

	paymentDoc := bson.M{
		"transaction_uuid": transactionUUID,
		"order_uuid":       payment.OrderUuid,
		"user_uuid":        payment.UserUuid,
		"payment_method":   int32(payment.PaymentMethod),
		"status":           "completed",
		"created_at":       now,
		"updated_at":       now,
	}

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_INITDB_DATABASE")
	if databaseName == "" {
		databaseName = "payments" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(paymentsCollectionName).InsertOne(ctx, paymentDoc)
	if err != nil {
		return "", err
	}

	return transactionUUID, nil
}

// GetTestPaymentData — возвращает тестовые данные платежа
func (env *TestEnvironment) GetTestPaymentData() *paymentv1.Payment {
	return &paymentv1.Payment{
		OrderUuid:     gofakeit.UUID(),
		UserUuid:      gofakeit.UUID(),
		PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
	}
}

// ClearPaymentsCollection — удаляет все записи из коллекции payments
func (env *TestEnvironment) ClearPaymentsCollection(ctx context.Context) error {
	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_INITDB_DATABASE")
	if databaseName == "" {
		databaseName = "payments" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(paymentsCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
