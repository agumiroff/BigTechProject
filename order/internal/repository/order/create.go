package order

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func (r *repository) CreateOrder(ctx context.Context, order *model.OrderRow, parts []string) (string, error) {
	const orderQuery = `
		INSERT INTO orders (
			order_uuid,
			user_uuid,
			total_price,
			status,
			created_at
		) VALUES (
			$1, $2, $3, $4, NOW()
		)`

	const orderPartsQuery = `
		INSERT INTO order_parts (
			order_uuid,
			part_uuid
		) SELECT $1, unnest($2::uuid[])`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	order.OrderUUID = uuid.New().String()
	_, err = tx.ExecContext(ctx, orderQuery,
		order.OrderUUID,
		order.UserUUID,
		order.TotalPrice,
		order.Status,
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return "", apperrors.ErrAlreadyExists
		}
		return "", fmt.Errorf("failed to create order: %w", err)
	}

	_, err = tx.ExecContext(ctx, orderPartsQuery,
		order.OrderUUID,
		pq.Array(parts),
	)
	if err != nil {
		return "", fmt.Errorf("failed to bulk insert order part: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return order.OrderUUID, nil
}
