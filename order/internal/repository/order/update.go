package order

import (
	"context"
	"database/sql"
	"fmt"

	repomodel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

const (
	getOrderStatusForUpdateQuery = `
		SELECT status FROM orders WHERE order_uuid = $1
	`
	updateOrderQuery = `
		UPDATE orders 
		SET user_uuid = $1,
			total_price = $2,
			status = $3,
			updated_at = NOW(),
			transaction_uuid = $4,
			payment_method = $5
		WHERE order_uuid = $6
	`
)

func (r *repository) UpdateOrder(ctx context.Context, m *repomodel.OrderRow) error {
	if m == nil {
		return apperrors.ErrInvalidRequest
	}

	if m.OrderUUID == "" {
		return apperrors.ErrInvalidRequest
	}

	// First check if order exists and get its current status
	var currentStatus string
	err := r.db.QueryRowContext(ctx, getOrderStatusForUpdateQuery, m.OrderUUID).Scan(&currentStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return apperrors.ErrNotFound
		}
		return fmt.Errorf("failed to get order status: %w", err)
	}

	if currentStatus == string(repomodel.OrderStatusPAID) &&
		m.Status != string(repomodel.OrderStatusCANCELLED) {
		return apperrors.ErrForbidden
	}

	// Update order
	result, err := r.db.ExecContext(ctx, updateOrderQuery,
		m.UserUUID,
		m.TotalPrice,
		m.Status,
		m.TransactionUUID.String,
		m.PaymentMethod.String,
		m.OrderUUID)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}
