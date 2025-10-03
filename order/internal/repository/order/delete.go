package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func (r *repository) DeleteOrder(ctx context.Context, uuid string) error {
	const (
		getStatusQuery = `SELECT status FROM orders WHERE order_uuid = $1`
		deleteQuery    = `DELETE FROM orders WHERE order_uuid = $1`
	)

	if uuid == "" {
		return apperrors.ErrInvalidRequest
	}

	// First check if order exists and get its current status
	var status model.OrderStatus
	err := r.db.QueryRowContext(ctx, getStatusQuery, uuid).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperrors.ErrNotFound
		}
		return fmt.Errorf("failed to check order status: %w", err)
	}

	if status == model.OrderStatusCANCELLED {
		return apperrors.ErrForbidden
	}

	// Delete order
	result, err := r.db.ExecContext(ctx, deleteQuery, uuid)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
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
