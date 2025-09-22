package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	"github.com/agumiroff/BigTechProject/order/v1/internal/repository/order"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func newTestOrder() *model.Order {
	return &model.Order{
		OrderUUID:  "test-order-uuid",
		UserUUID:   "test-user-uuid",
		PartUUIDs:  []string{"part1", "part2"},
		TotalPrice: 100.50,
		Status:     model.OrderStatusPENDINGPAYMENT,
	}
}

func TestCreateOrder_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := order.NewRepository()
	testOrder := newTestOrder()

	// Act
	resp, err := repo.CreateOrder(ctx, testOrder)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, testOrder.OrderUUID, resp.OrderUUID)
	assert.Equal(t, testOrder.TotalPrice, resp.TotalPrice)

	// Verify storage
	stored, err := repo.Get(ctx, testOrder.OrderUUID)
	require.NoError(t, err)
	require.NotNil(t, stored)
	assert.Equal(t, testOrder.UserUUID, stored.UserUUID)
	assert.Equal(t, testOrder.PartUUIDs, stored.PartUUIDs)
	assert.Equal(t, model.OrderStatus(testOrder.Status), model.OrderStatus(stored.Status))
}

func TestCreateOrder_Overwrite(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := order.NewRepository()
	existingOrder := &model.Order{
		OrderUUID:  "same-uuid",
		UserUUID:   "old-user",
		PartUUIDs:  []string{"old-part"},
		TotalPrice: 50.00,
		Status:     model.OrderStatusPENDINGPAYMENT,
	}
	_, err := repo.CreateOrder(ctx, existingOrder)
	require.NoError(t, err)

	newOrder := &model.Order{
		OrderUUID:  "same-uuid",
		UserUUID:   "new-user",
		PartUUIDs:  []string{"new-part"},
		TotalPrice: 150.00,
		Status:     model.OrderStatusPAID,
	}

	// Act
	resp, err := repo.CreateOrder(ctx, newOrder)

	// Assert
	require.ErrorIs(t, err, apperrors.ErrAlreadyExists)
	assert.Nil(t, resp)

	// Verify storage unchanged
	stored, err := repo.Get(ctx, newOrder.OrderUUID)
	require.NoError(t, err)
	require.NotNil(t, stored)
	assert.Equal(t, existingOrder.UserUUID, stored.UserUUID)
	assert.Equal(t, existingOrder.PartUUIDs, stored.PartUUIDs)
	assert.Equal(t, model.OrderStatus(existingOrder.Status), model.OrderStatus(stored.Status))
}

func TestCreateOrder_NilOrder(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := order.NewRepository()

	// Act
	resp, err := repo.CreateOrder(ctx, nil)

	// Assert
	assert.ErrorIs(t, err, apperrors.ErrInvalidRequest)
	assert.Nil(t, resp)
}

func TestCreateOrder_EmptyOrderUUID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := order.NewRepository()
	invalidOrder := &model.Order{
		UserUUID:   "test-user",
		PartUUIDs:  []string{"part1"},
		TotalPrice: 100,
		Status:     model.OrderStatusPENDINGPAYMENT,
	}

	// Act
	resp, err := repo.CreateOrder(ctx, invalidOrder)

	// Assert
	assert.ErrorIs(t, err, apperrors.ErrInvalidRequest)
	assert.Nil(t, resp)

	// Verify storage wasn't modified
	stored, err := repo.Get(ctx, "")
	assert.ErrorIs(t, err, apperrors.ErrInvalidRequest)
	assert.Nil(t, stored)
}

func TestCreateOrder_EmptyUserUUID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := order.NewRepository()
	invalidOrder := &model.Order{
		OrderUUID:  "test-uuid",
		PartUUIDs:  []string{"part1"},
		TotalPrice: 100,
		Status:     model.OrderStatusPENDINGPAYMENT,
	}

	// Act
	resp, err := repo.CreateOrder(ctx, invalidOrder)

	// Assert
	assert.ErrorIs(t, err, apperrors.ErrInvalidRequest)
	assert.Nil(t, resp)

	// Verify storage wasn't modified
	stored, err := repo.Get(ctx, invalidOrder.OrderUUID)
	assert.ErrorIs(t, err, apperrors.ErrNotFound)
	assert.Nil(t, stored)
}

func TestCreateOrder_EmptyPartUUIDs(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := order.NewRepository()
	invalidOrder := &model.Order{
		OrderUUID:  "test-uuid",
		UserUUID:   "test-user",
		TotalPrice: 100,
		Status:     model.OrderStatusPENDINGPAYMENT,
	}

	// Act
	resp, err := repo.CreateOrder(ctx, invalidOrder)

	// Assert
	assert.ErrorIs(t, err, apperrors.ErrInvalidRequest)
	assert.Nil(t, resp)

	// Verify storage wasn't modified
	stored, err := repo.Get(ctx, invalidOrder.OrderUUID)
	assert.ErrorIs(t, err, apperrors.ErrNotFound)
	assert.Nil(t, stored)
}
