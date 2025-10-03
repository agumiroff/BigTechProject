package unit

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	repomodel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func TestUpdateOrder_Success(t *testing.T) {
	// Arrange
	repo := NewInmemoryRepo()
	ctx := context.Background()

	orderUUID := "test-order-uuid"
	userUUID := "test-user-uuid"
	parts := []string{"part1"}
	totalPrice := 100.0

	// Create initial order
	initialOrder := &repomodel.OrderRow{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		TotalPrice: totalPrice,
		Status:     string(model.OrderStatusPENDINGPAYMENT),
	}

	_, err := repo.CreateOrder(ctx, initialOrder, parts)
	require.NoError(t, err)

	// Create updated order with payment info
	updatedOrder := &repomodel.OrderRow{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		TotalPrice:      totalPrice,
		Status:          string(model.OrderStatusPAID),
		TransactionUUID: sql.NullString{String: "transaction-123", Valid: true},
		PaymentMethod:   sql.NullString{String: string(model.PaymentMethodCARD), Valid: true},
	}

	// Act
	err = repo.UpdateOrder(ctx, updatedOrder)

	// Assert
	require.NoError(t, err)

	// Verify the order was updated
	storedOrder, _, err := repo.GetOrder(ctx, orderUUID)
	require.NoError(t, err)
	assert.Equal(t, updatedOrder.Status, storedOrder.Status)
	assert.True(t, storedOrder.TransactionUUID.Valid)
	assert.Equal(t, "transaction-123", storedOrder.TransactionUUID.String)
	assert.True(t, storedOrder.PaymentMethod.Valid)
	assert.Equal(t, string(model.PaymentMethodCARD), storedOrder.PaymentMethod.String)
}

func TestUpdateOrder_OrderNotFound(t *testing.T) {
	// Arrange
	repo := NewInmemoryRepo()
	ctx := context.Background()

	order := &repomodel.OrderRow{
		OrderUUID:  "nonexistent-uuid",
		UserUUID:   "test-user",
		TotalPrice: 100.0,
		Status:     string(model.OrderStatusPAID),
	}

	// Act
	err := repo.UpdateOrder(ctx, order)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrNotFound)
}

func TestUpdateOrder_NilOrder(t *testing.T) {
	// Arrange
	repo := NewInmemoryRepo()
	ctx := context.Background()

	// Act
	err := repo.UpdateOrder(ctx, nil)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrInvalidRequest)
}

func TestUpdateOrder_EmptyUUID(t *testing.T) {
	// Arrange
	repo := NewInmemoryRepo()
	ctx := context.Background()

	order := &repomodel.OrderRow{
		OrderUUID:  "", // Empty UUID
		UserUUID:   "test-user",
		TotalPrice: 100.0,
		Status:     string(model.OrderStatusPAID),
	}

	// Act
	err := repo.UpdateOrder(ctx, order)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrInvalidRequest)
}

func TestUpdateOrder_CantUpdateCancelledOrder(t *testing.T) {
	// Arrange
	repo := NewInmemoryRepo()
	ctx := context.Background()

	orderUUID := "test-order-uuid"
	userUUID := "test-user-uuid"
	parts := []string{"part1"}

	// Create initial order
	initialOrder := &repomodel.OrderRow{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		TotalPrice: 100.0,
		Status:     string(model.OrderStatusCANCELLED),
	}

	_, err := repo.CreateOrder(ctx, initialOrder, parts)
	require.NoError(t, err)

	// Try to update cancelled order
	updatedOrder := &repomodel.OrderRow{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		TotalPrice: 150.0, // Changed price
		Status:     string(model.OrderStatusPENDINGPAYMENT),
	}

	// Act
	err = repo.UpdateOrder(ctx, updatedOrder)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrForbidden)

	// Verify the order wasn't updated
	storedOrder, _, err := repo.GetOrder(ctx, orderUUID)
	require.NoError(t, err)
	assert.Equal(t, string(model.OrderStatusCANCELLED), storedOrder.Status)
	assert.Equal(t, float64(100.0), storedOrder.TotalPrice)
}

func TestUpdateOrder_CantChangeStatusFromPaid(t *testing.T) {
	// Arrange
	repo := NewInmemoryRepo()
	ctx := context.Background()

	orderUUID := "test-order-uuid"
	userUUID := "test-user-uuid"
	parts := []string{"part1"}

	// Create initial paid order
	initialOrder := &repomodel.OrderRow{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		TotalPrice:      100.0,
		Status:          string(model.OrderStatusPAID),
		TransactionUUID: sql.NullString{String: "transaction-123", Valid: true},
		PaymentMethod:   sql.NullString{String: string(model.PaymentMethodCARD), Valid: true},
	}

	_, err := repo.CreateOrder(ctx, initialOrder, parts)
	require.NoError(t, err)

	// Try to change status from PAID to PENDING
	updatedOrder := &repomodel.OrderRow{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		TotalPrice:      100.0,
		Status:          string(model.OrderStatusPENDINGPAYMENT),
		TransactionUUID: sql.NullString{String: "transaction-123", Valid: true},
		PaymentMethod:   sql.NullString{String: string(model.PaymentMethodCARD), Valid: true},
	}

	// Act
	err = repo.UpdateOrder(ctx, updatedOrder)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrForbidden)

	// Verify the order wasn't updated
	storedOrder, _, err := repo.GetOrder(ctx, orderUUID)
	require.NoError(t, err)
	assert.Equal(t, string(model.OrderStatusPAID), storedOrder.Status)
}

func TestUpdateOrder_CanChangePaidToCancelled(t *testing.T) {
	// Arrange
	repo := NewInmemoryRepo()
	ctx := context.Background()

	orderUUID := "test-order-uuid"
	userUUID := "test-user-uuid"
	parts := []string{"part1"}

	// Create initial paid order
	initialOrder := &repomodel.OrderRow{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		TotalPrice:      100.0,
		Status:          string(model.OrderStatusPAID),
		TransactionUUID: sql.NullString{String: "transaction-123", Valid: true},
		PaymentMethod:   sql.NullString{String: string(model.PaymentMethodCARD), Valid: true},
	}

	_, err := repo.CreateOrder(ctx, initialOrder, parts)
	require.NoError(t, err)

	// Try to change status from PAID to CANCELLED (allowed)
	updatedOrder := &repomodel.OrderRow{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		TotalPrice:      100.0,
		Status:          string(model.OrderStatusCANCELLED),
		TransactionUUID: sql.NullString{String: "transaction-123", Valid: true},
		PaymentMethod:   sql.NullString{String: string(model.PaymentMethodCARD), Valid: true},
	}

	// Act
	err = repo.UpdateOrder(ctx, updatedOrder)

	// Assert
	require.NoError(t, err)

	// Verify the order was updated to cancelled
	storedOrder, _, err := repo.GetOrder(ctx, orderUUID)
	require.NoError(t, err)
	assert.Equal(t, string(model.OrderStatusCANCELLED), storedOrder.Status)
}
