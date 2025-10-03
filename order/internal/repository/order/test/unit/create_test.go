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

func TestCreateOrder_Success(t *testing.T) {
	// Arrange
	repo := NewInmemoryRepo()
	ctx := context.Background()

	orderUUID := "test-order-uuid"
	userUUID := "test-user-uuid"
	parts := []string{"part1", "part2"}
	totalPrice := 100.0

	order := &repomodel.OrderRow{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		TotalPrice: totalPrice,
		Status:     string(model.OrderStatusPENDINGPAYMENT),
	}

	// Act
	resp, err := repo.CreateOrder(ctx, order, parts)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, orderUUID, resp.OrderUUID)
	assert.Equal(t, totalPrice, resp.TotalPrice)

	// Verify order was stored
	storedOrder, storedParts, err := repo.GetOrder(ctx, orderUUID)
	require.NoError(t, err)
	require.NotNil(t, storedOrder)
	assert.Equal(t, order.OrderUUID, storedOrder.OrderUUID)
	assert.Equal(t, order.UserUUID, storedOrder.UserUUID)
	assert.Equal(t, order.TotalPrice, storedOrder.TotalPrice)
	assert.Equal(t, order.Status, storedOrder.Status)
	assert.Equal(t, parts, storedParts)
}

func TestCreateOrder_AlreadyExists(t *testing.T) {
	// Arrange
	repo := NewInmemoryRepo()
	ctx := context.Background()

	orderUUID := "test-order-uuid"
	userUUID := "test-user-uuid"
	parts := []string{"part1", "part2"}
	totalPrice := 100.0

	order := &repomodel.OrderRow{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		TotalPrice: totalPrice,
		Status:     string(model.OrderStatusPENDINGPAYMENT),
	}

	// Create the first order
	_, err := repo.CreateOrder(ctx, order, parts)
	require.NoError(t, err)

	// Act - Try to create order with the same UUID
	resp, err := repo.CreateOrder(ctx, order, parts)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrAlreadyExists)
	assert.Nil(t, resp)
}

func TestCreateOrder_InvalidRequest_NilOrder(t *testing.T) {
	// Arrange
	repo := NewInmemoryRepo()
	ctx := context.Background()

	// Act
	resp, err := repo.CreateOrder(ctx, nil, []string{"part1"})

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrInvalidRequest)
	assert.Nil(t, resp)
}

func TestCreateOrder_InvalidRequest_EmptyUUID(t *testing.T) {
	// Arrange
	repo := NewInmemoryRepo()
	ctx := context.Background()

	order := &repomodel.OrderRow{
		OrderUUID:  "", // Empty UUID
		UserUUID:   "test-user-uuid",
		TotalPrice: 100.0,
		Status:     string(model.OrderStatusPENDINGPAYMENT),
	}

	// Act
	resp, err := repo.CreateOrder(ctx, order, []string{"part1"})

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrInvalidRequest)
	assert.Nil(t, resp)
}

func TestCreateOrder_InvalidRequest_EmptyUserUUID(t *testing.T) {
	// Arrange
	repo := NewInmemoryRepo()
	ctx := context.Background()

	order := &repomodel.OrderRow{
		OrderUUID:  "test-order-uuid",
		UserUUID:   "", // Empty UserUUID
		TotalPrice: 100.0,
		Status:     string(model.OrderStatusPENDINGPAYMENT),
	}

	// Act
	resp, err := repo.CreateOrder(ctx, order, []string{"part1"})

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrInvalidRequest)
	assert.Nil(t, resp)
}

func TestCreateOrder_InvalidRequest_EmptyParts(t *testing.T) {
	// Arrange
	repo := NewInmemoryRepo()
	ctx := context.Background()

	order := &repomodel.OrderRow{
		OrderUUID:  "test-order-uuid",
		UserUUID:   "test-user-uuid",
		TotalPrice: 100.0,
		Status:     string(model.OrderStatusPENDINGPAYMENT),
	}

	// Act
	resp, err := repo.CreateOrder(ctx, order, []string{}) // Empty parts

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrInvalidRequest)
	assert.Nil(t, resp)
}

func TestCreateOrder_WithNullFields(t *testing.T) {
	// Arrange
	repo := NewInmemoryRepo()
	ctx := context.Background()

	orderUUID := "test-order-uuid"
	userUUID := "test-user-uuid"
	parts := []string{"part1"}
	totalPrice := 100.0

	order := &repomodel.OrderRow{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		TotalPrice: totalPrice,
		Status:     string(model.OrderStatusPENDINGPAYMENT),
		// SQL null values
		TransactionUUID: sql.NullString{Valid: false},
		PaymentMethod:   sql.NullString{Valid: false},
	}

	// Act
	resp, err := repo.CreateOrder(ctx, order, parts)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, orderUUID, resp.OrderUUID)
	assert.Equal(t, totalPrice, resp.TotalPrice)

	// Verify order was stored with null fields
	storedOrder, _, err := repo.GetOrder(ctx, orderUUID)
	require.NoError(t, err)
	require.NotNil(t, storedOrder)
	assert.Equal(t, false, storedOrder.TransactionUUID.Valid)
	assert.Equal(t, false, storedOrder.PaymentMethod.Valid)
}
