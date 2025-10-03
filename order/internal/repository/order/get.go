package order

import (
	"context"
	"database/sql"
	"fmt"

	repomodel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func (r *repository) GetOrder(ctx context.Context, uuid string) (order *repomodel.OrderRow, parts []string, err error) {
	if uuid == "" {
		return nil, nil, apperrors.ErrInvalidRequest
	}

	const orderQuery = `
		SELECT 
			order_uuid,
			user_uuid,
			total_price,
			status,
			payment_method,
			transaction_uuid
		FROM orders 
		WHERE order_uuid = $1`

	const partsQuery = `
		SELECT 
			part_uuid
		FROM order_parts 
		WHERE order_uuid = $1`

	order = &repomodel.OrderRow{}
	err = r.db.QueryRowContext(ctx, orderQuery, uuid).Scan(
		&order.OrderUUID,
		&order.UserUUID,
		&order.TotalPrice,
		&order.Status,
		&order.PaymentMethod,
		&order.TransactionUUID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, apperrors.ErrNotFound
		}
		return nil, nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Get order parts
	rows, err := r.db.QueryContext(ctx, partsQuery, uuid)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get order parts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var partUUID string
		if err = rows.Scan(&partUUID); err != nil {
			return nil, nil, fmt.Errorf("failed to scan part UUID: %w", err)
		}
		parts = append(parts, partUUID)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error iterating order parts: %w", err)
	}

	return order, parts, nil
}
