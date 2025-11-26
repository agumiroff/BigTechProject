package integration

import (
	"context"
	"fmt"

	"github.com/brianvoe/gofakeit/v6"

	orderv1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
)

// CreateOrder — создает новый заказ
func (env *TestEnvironment) CreateOrder(ctx context.Context) orderv1.Order {
	order := orderv1.Order{
		OrderUUID: gofakeit.UUID(),
		UserUUID:  gofakeit.UUID(),
		PartUuids: []string{
			gofakeit.UUID(),
			gofakeit.UUID(),
		},
		TotalPrice: 149.99,
		TransactionUUID: orderv1.OptNilString{
			Value: gofakeit.UUID(),
		},
		PaymentMethod: orderv1.OptPaymentMethod{
			Value: orderv1.PaymentMethodCARD,
		},
		Status: orderv1.OrderStatusPENDINGPAYMENT,
	}

	return order
}

// ClearOrderCollection — удаляет все записи из таблиц orders и order_parts
func (env *TestEnvironment) ClearOrderCollection(ctx context.Context) error {
	pool := env.Postgres.Pool()

	// Очищаем таблицы в правильном порядке (CASCADE автоматически очистит связанные таблицы)
	queries := []string{
		"TRUNCATE TABLE order_parts CASCADE",
		"TRUNCATE TABLE orders CASCADE",
	}

	for _, query := range queries {
		if _, err := pool.Exec(ctx, query); err != nil {
			return fmt.Errorf("failed to truncate table: %w", err)
		}
	}

	return nil
}
