package order

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func (r *repository) CreateOrder(ctx context.Context, req *model.Order) (*model.CreateOrderResponse, error) {
	const query = `
		INSERT INTO orders (
			order_uuid,
			user_uuid,
			part_uuids,
			total_price,
			status,
			created_at
		) VALUES (
			$1, $2, $3, $4, $5, NOW()
		) RETURNING order_uuid, total_price`

	if req == nil {
		return nil, apperrors.ErrInvalidRequest
	}

	if req.OrderUUID == "" {
		return nil, apperrors.ErrInvalidRequest
	}

	if req.UserUUID == "" || len(req.PartUUIDs) == 0 {
		return nil, apperrors.ErrInvalidRequest
	}

	// Convert part UUIDs to JSON
	partUUIDsJSON, err := json.Marshal(req.PartUUIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal part UUIDs: %w", err)
	}

	// Insert order into database
	var orderUUID string
	var totalPrice float64

	err = r.db.QueryRowContext(ctx, query,
		req.OrderUUID,
		req.UserUUID,
		string(partUUIDsJSON),
		req.TotalPrice,
		req.Status,
	).Scan(&orderUUID, &totalPrice)

	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"orders_pkey\"" {
			return nil, apperrors.ErrAlreadyExists
		}
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return &model.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}
