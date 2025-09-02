package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	rModel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/order/v1/internal/repository/order"
)

func TestDeleteOrder_Success(t *testing.T) {
	// Arrange
	repo := order.NewRepository()
	testOrder := newTestOrder()
	repo.CreateOrder(testOrder)

	// Verify order exists
	_, err := repo.Get(testOrder.OrderUUID)
	require.NoError(t, err)

	// Act
	err = repo.DeleteOrder(testOrder.OrderUUID)

	// Assert
	assert.NoError(t, err)

	// Verify deletion
	stored, err := repo.Get(testOrder.OrderUUID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, rModel.ErrOrderNotFound)
	assert.Nil(t, stored)
}

func TestDeleteOrder_NotFound(t *testing.T) {
	// Arrange
	repo := order.NewRepository()

	// Act
	err := repo.DeleteOrder("non-existent-uuid")

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, rModel.ErrOrderNotFound)
}

func TestDeleteOrder_EmptyUUID(t *testing.T) {
	// Arrange
	repo := order.NewRepository()

	// Act
	err := repo.DeleteOrder("")

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, rModel.ErrInvalidOrderUUID)
}

func TestDeleteOrder_CancelledOrder(t *testing.T) {
	// Arrange
	repo := order.NewRepository()
	testOrder := newTestOrder()
	testOrder.Status = model.OrderStatusCANCELLED
	repo.CreateOrder(testOrder)

	// Act
	err := repo.DeleteOrder(testOrder.OrderUUID)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, rModel.ErrOrderAlreadyCancelled)

	// Verify order still exists
	stored, err := repo.Get(testOrder.OrderUUID)
	require.NoError(t, err)
	assert.NotNil(t, stored)
	assert.Equal(t, rModel.OrderStatusCANCELLED, stored.Status)
}
