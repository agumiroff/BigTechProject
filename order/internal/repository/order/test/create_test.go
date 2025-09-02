package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	rModel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/order/v1/internal/repository/order"
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
	repo := order.NewRepository()
	testOrder := newTestOrder()

	// Act
	resp := repo.CreateOrder(testOrder)

	// Assert
	require.NotNil(t, resp)
	assert.Equal(t, testOrder.OrderUUID, resp.OrderUUID)
	assert.Equal(t, testOrder.TotalPrice, resp.TotalPrice)

	// Verify storage
	stored, err := repo.Get(testOrder.OrderUUID)
	require.NoError(t, err)
	require.NotNil(t, stored)
	assert.Equal(t, testOrder.UserUUID, stored.UserUUID)
	assert.Equal(t, testOrder.PartUUIDs, stored.PartUUIDs)
	assert.Equal(t, rModel.OrderStatus(testOrder.Status), stored.Status)
}

func TestCreateOrder_Overwrite(t *testing.T) {
	// Arrange
	repo := order.NewRepository()
	existingOrder := &model.Order{
		OrderUUID:  "same-uuid",
		UserUUID:   "old-user",
		PartUUIDs:  []string{"old-part"},
		TotalPrice: 50.00,
		Status:     model.OrderStatusPENDINGPAYMENT,
	}
	repo.CreateOrder(existingOrder)

	newOrder := &model.Order{
		OrderUUID:  "same-uuid",
		UserUUID:   "new-user",
		PartUUIDs:  []string{"new-part"},
		TotalPrice: 150.00,
		Status:     model.OrderStatusPAID,
	}

	// Act
	resp := repo.CreateOrder(newOrder)

	// Assert
	require.NotNil(t, resp)
	assert.Equal(t, newOrder.OrderUUID, resp.OrderUUID)
	assert.Equal(t, newOrder.TotalPrice, resp.TotalPrice)

	// Verify storage
	stored, err := repo.Get(newOrder.OrderUUID)
	require.NoError(t, err)
	require.NotNil(t, stored)
	assert.Equal(t, newOrder.UserUUID, stored.UserUUID)
	assert.Equal(t, newOrder.PartUUIDs, stored.PartUUIDs)
	assert.Equal(t, rModel.OrderStatus(newOrder.Status), stored.Status)
}

func TestCreateOrder_NilOrder(t *testing.T) {
	// Arrange
	repo := order.NewRepository()

	// Act
	resp := repo.CreateOrder(nil)

	// Assert
	assert.Nil(t, resp)
}

func TestCreateOrder_EmptyOrderUUID(t *testing.T) {
	// Arrange
	repo := order.NewRepository()
	invalidOrder := &model.Order{
		UserUUID:   "test-user",
		PartUUIDs:  []string{"part1"},
		TotalPrice: 100,
		Status:     model.OrderStatusPENDINGPAYMENT,
	}

	// Act
	resp := repo.CreateOrder(invalidOrder)

	// Assert
	assert.Nil(t, resp)

	// Verify storage wasn't modified
	stored, err := repo.Get("")
	assert.Error(t, err)
	assert.ErrorIs(t, err, rModel.ErrInvalidOrderUUID)
	assert.Nil(t, stored)
}

func TestCreateOrder_EmptyUserUUID(t *testing.T) {
	// Arrange
	repo := order.NewRepository()
	invalidOrder := &model.Order{
		OrderUUID:  "test-uuid",
		PartUUIDs:  []string{"part1"},
		TotalPrice: 100,
		Status:     model.OrderStatusPENDINGPAYMENT,
	}

	// Act
	resp := repo.CreateOrder(invalidOrder)

	// Assert
	assert.Nil(t, resp)

	// Verify storage wasn't modified
	stored, err := repo.Get(invalidOrder.OrderUUID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, rModel.ErrOrderNotFound)
	assert.Nil(t, stored)
}

func TestCreateOrder_EmptyPartUUIDs(t *testing.T) {
	// Arrange
	repo := order.NewRepository()
	invalidOrder := &model.Order{
		OrderUUID:  "test-uuid",
		UserUUID:   "test-user",
		TotalPrice: 100,
		Status:     model.OrderStatusPENDINGPAYMENT,
	}

	// Act
	resp := repo.CreateOrder(invalidOrder)

	// Assert
	assert.Nil(t, resp)

	// Verify storage wasn't modified
	stored, err := repo.Get(invalidOrder.OrderUUID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, rModel.ErrOrderNotFound)
	assert.Nil(t, stored)
}
