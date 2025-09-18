package order

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	repomodel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func (r *repository) Get(ctx context.Context, uuid string) (*repomodel.Order, error) {
	const query = `
		SELECT 
			order_uuid,
			user_uuid,
			part_uuids,
			total_price,
			status,
			payment_method,
			transaction_uuid
		FROM orders 
		WHERE order_uuid = $1`

	if uuid == "" {
		return nil, apperrors.ErrInvalidRequest
	}

	var order repomodel.Order
	var partUUIDsJSON string

	var paymentMethod sql.NullString
	var transactionUUID sql.NullString
	err := r.db.QueryRowContext(ctx, query, uuid).Scan(
		&order.OrderUUID,
		&order.UserUUID,
		&partUUIDsJSON,
		&order.TotalPrice,
		&order.Status,
		&paymentMethod,
		&transactionUUID,
	)

	if paymentMethod.Valid {
		order.PaymentMethod = repomodel.PaymentMethod(paymentMethod.String)
	} else {
		order.PaymentMethod = repomodel.PaymentMethodUNKNOWN
	}

	if transactionUUID.Valid {
		order.TransactionUUID = ""
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	var partUUIDs []string
	if err := json.Unmarshal([]byte(partUUIDsJSON), &partUUIDs); err != nil {
		return nil, fmt.Errorf("failed to parse part UUIDs: %w", err)
	}
	order.PartUUIDs = partUUIDs

	return &order, nil
}
