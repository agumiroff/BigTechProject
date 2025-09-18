package order

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func (r *repository) UpdateOrder(ctx context.Context, m *model.Order) error {
	if m == nil {
		return apperrors.ErrInvalidRequest
	}

	if m.OrderUUID == "" {
		return apperrors.ErrInvalidRequest
	}

	// First check if order exists and get its current status
	var currentStatus model.OrderStatus
	err := r.db.QueryRowContext(ctx, `
		SELECT status FROM orders WHERE order_uuid = $1
	`, m.OrderUUID).Scan(&currentStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return apperrors.ErrNotFound
		}
		return fmt.Errorf("failed to get order status: %w", err)
	}

	if currentStatus == model.OrderStatusPAID &&
		m.Status != model.OrderStatusCANCELLED {
		return apperrors.ErrForbidden
	}

	// Convert part UUIDs to JSON
	partUUIDsJSON, err := json.Marshal(m.PartUUIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal part UUIDs: %w", err)
	}

	// Update order
	result, err := r.db.ExecContext(ctx, `
		UPDATE orders 
		SET user_uuid = $1,
			part_uuids = $2,
			total_price = $3,
			status = $4,
			updated_at = NOW()
		WHERE order_uuid = $5
	`, m.UserUUID, string(partUUIDsJSON), m.TotalPrice, m.Status, m.OrderUUID)
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
