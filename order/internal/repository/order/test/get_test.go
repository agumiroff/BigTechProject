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

func TestGet_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := order.NewRepository()
	testOrder := newTestOrder()
	_, err := repo.CreateOrder(ctx, testOrder)
	require.NoError(t, err)

	// Act
	stored, err := repo.Get(ctx, testOrder.OrderUUID)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, stored)
	assert.Equal(t, testOrder.UserUUID, stored.UserUUID)
	assert.Equal(t, testOrder.PartUUIDs, stored.PartUUIDs)
	assert.Equal(t, model.OrderStatus(testOrder.Status), model.OrderStatus(stored.Status))
}

func TestGet_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := order.NewRepository()

	// Act
	stored, err := repo.Get(ctx, "non-existent-uuid")

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrNotFound)
	assert.Nil(t, stored)
}

func TestGet_EmptyUUID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := order.NewRepository()

	// Act
	stored, err := repo.Get(ctx, "")

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrInvalidRequest)
	assert.Nil(t, stored)
}
