package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	"github.com/agumiroff/BigTechProject/order/v1/internal/repository/order/inmemory"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func TestDeleteOrder_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := inmemory.NewInMemoryOrderRepository()
	testOrder := newTestOrder()
	_, err := repo.CreateOrder(ctx, testOrder)
	require.NoError(t, err)

	// Verify order exists
	_, err = repo.Get(ctx, testOrder.OrderUUID)
	require.NoError(t, err)

	// Act
	err = repo.DeleteOrder(ctx, testOrder.OrderUUID)

	// Assert
	assert.NoError(t, err)

	// Verify deletion
	stored, err := repo.Get(ctx, testOrder.OrderUUID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrNotFound)
	assert.Nil(t, stored)
}

func TestDeleteOrder_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := inmemory.NewInMemoryOrderRepository()

	// Act
	err := repo.DeleteOrder(ctx, "non-existent-uuid")

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrNotFound)
}

func TestDeleteOrder_EmptyUUID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := inmemory.NewInMemoryOrderRepository()

	// Act
	err := repo.DeleteOrder(ctx, "")

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrInvalidRequest)
}

func TestDeleteOrder_CancelledOrder(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := inmemory.NewInMemoryOrderRepository()
	testOrder := newTestOrder()
	testOrder.Status = model.OrderStatusCANCELLED
	_, err := repo.CreateOrder(ctx, testOrder)
	require.NoError(t, err)

	// Act
	err = repo.DeleteOrder(ctx, testOrder.OrderUUID)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrForbidden)

	// Verify order still exists
	stored, err := repo.Get(ctx, testOrder.OrderUUID)
	require.NoError(t, err)
	assert.NotNil(t, stored)
	assert.Equal(t, model.OrderStatusCANCELLED, model.OrderStatus(stored.Status))
}
