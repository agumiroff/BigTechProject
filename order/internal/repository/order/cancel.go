package order

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func (r *repository) CancelOrder(ctx context.Context, uuid string) error {
	if uuid == "" {
		return apperrors.ErrInvalidRequest
	}

	// First check if the order exists and get its current status
	var currentStatus model.OrderStatus
	err := r.db.QueryRowContext(ctx, `
		SELECT status FROM orders WHERE order_uuid = $1
	`, uuid).Scan(&currentStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return apperrors.ErrNotFound
		}
		return fmt.Errorf("failed to get order status: %w", err)
	}

	if currentStatus == model.OrderStatusCANCELLED {
		return apperrors.ErrForbidden
	}

	if currentStatus == model.OrderStatusPAID {
		return apperrors.ErrForbidden
	}

	// Update order status to cancelled
	result, err := r.db.ExecContext(ctx, `
		UPDATE orders 
		SET status = $1, updated_at = NOW()
		WHERE order_uuid = $2
	`, model.OrderStatusCANCELLED, uuid)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
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
