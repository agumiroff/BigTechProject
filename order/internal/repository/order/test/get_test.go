package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	rModel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/order/v1/internal/repository/order"
)

func TestGet_Success(t *testing.T) {
	// Arrange
	repo := order.NewRepository()
	testOrder := newTestOrder()
	repo.CreateOrder(testOrder)

	// Act
	stored, err := repo.Get(testOrder.OrderUUID)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, stored)
	assert.Equal(t, testOrder.UserUUID, stored.UserUUID)
	assert.Equal(t, testOrder.PartUUIDs, stored.PartUUIDs)
	assert.Equal(t, rModel.OrderStatus(testOrder.Status), stored.Status)
}

func TestGet_NotFound(t *testing.T) {
	// Arrange
	repo := order.NewRepository()

	// Act
	stored, err := repo.Get("non-existent-uuid")

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, rModel.ErrOrderNotFound)
	assert.Nil(t, stored)
}

func TestGet_EmptyUUID(t *testing.T) {
	// Arrange
	repo := order.NewRepository()

	// Act
	stored, err := repo.Get("")

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, rModel.ErrInvalidOrderUUID)
	assert.Nil(t, stored)
}
