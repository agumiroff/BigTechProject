package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func TestUpdateOrder_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := NewInmemoryRepo()
	testOrder := newTestOrder()
	_, err := repo.CreateOrder(ctx, testOrder)
	require.NoError(t, err)

	updatedOrder := &model.Order{
		OrderUUID:  testOrder.OrderUUID,
		UserUUID:   "new-user-uuid",
		PartUUIDs:  []string{"new-part-1", "new-part-2"},
		TotalPrice: 200.50,
		Status:     model.OrderStatusPAID,
	}

	// Act
	err = repo.UpdateOrder(ctx, updatedOrder)

	// Assert
	require.NoError(t, err)

	// Verify storage
	stored, err := repo.Get(ctx, testOrder.OrderUUID)
	require.NoError(t, err)
	require.NotNil(t, stored)
	assert.Equal(t, updatedOrder.UserUUID, stored.UserUUID)
	assert.Equal(t, updatedOrder.PartUUIDs, stored.PartUUIDs)
	assert.Equal(t, string(updatedOrder.Status), string(stored.Status))
}

func TestUpdateOrder_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := NewInmemoryRepo()
	testOrder := &model.Order{
		OrderUUID:  "non-existent-uuid",
		UserUUID:   "test-user",
		PartUUIDs:  []string{"part1"},
		TotalPrice: 100,
		Status:     model.OrderStatusPENDINGPAYMENT,
	}

	// Act
	err := repo.UpdateOrder(ctx, testOrder)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrNotFound)
}

func TestUpdateOrder_NilOrder(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := NewInmemoryRepo()

	// Act
	err := repo.UpdateOrder(ctx, nil)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrInvalidRequest)
}

func TestUpdateOrder_EmptyUUID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := NewInmemoryRepo()
	testOrder := &model.Order{
		UserUUID:   "test-user",
		PartUUIDs:  []string{"part1"},
		TotalPrice: 100,
		Status:     model.OrderStatusPENDINGPAYMENT,
	}

	// Act
	err := repo.UpdateOrder(ctx, testOrder)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrInvalidRequest)
}

func TestUpdateOrder_CancelledOrder(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := NewInmemoryRepo()
	testOrder := newTestOrder()
	testOrder.Status = model.OrderStatusCANCELLED
	_, err := repo.CreateOrder(ctx, testOrder)
	require.NoError(t, err)

	updateReq := &model.Order{
		OrderUUID:  testOrder.OrderUUID,
		UserUUID:   "new-user",
		PartUUIDs:  []string{"new-part"},
		TotalPrice: 200,
		Status:     model.OrderStatusPAID,
	}

	// Act
	err = repo.UpdateOrder(ctx, updateReq)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrForbidden)

	// Verify order wasn't modified
	stored, err := repo.Get(ctx, testOrder.OrderUUID)
	require.NoError(t, err)
	assert.Equal(t, string(model.OrderStatusCANCELLED), string(stored.Status))
}
