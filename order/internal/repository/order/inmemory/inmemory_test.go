package inmemory

import (
	"context"
	"testing"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	repomodel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryOrderRepository_CreateOrder(t *testing.T) {
	repo := NewInMemoryOrderRepository()
	ctx := context.Background()

	t.Run("successfully create order", func(t *testing.T) {
		order := &model.Order{
			OrderUUID:  uuid.New().String(),
			UserUUID:   uuid.New().String(),
			PartUUIDs:  []string{uuid.New().String(), uuid.New().String()},
			TotalPrice: 100.50,
			Status:     model.OrderStatusPENDINGPAYMENT,
		}

		resp, err := repo.CreateOrder(ctx, order)
		require.NoError(t, err)
		assert.Equal(t, order.OrderUUID, resp.OrderUUID)
		assert.Equal(t, order.TotalPrice, resp.TotalPrice)

		// Verify order was stored
		stored, err := repo.Get(ctx, order.OrderUUID)
		require.NoError(t, err)
		assert.Equal(t, order.OrderUUID, stored.OrderUUID)
		assert.Equal(t, order.UserUUID, stored.UserUUID)
		assert.Equal(t, order.PartUUIDs, stored.PartUUIDs)
		assert.Equal(t, order.TotalPrice, stored.TotalPrice)
		assert.Equal(t, repomodel.OrderStatus(order.Status), stored.Status)
	})

	t.Run("fail on duplicate order", func(t *testing.T) {
		order := &model.Order{
			OrderUUID:  uuid.New().String(),
			UserUUID:   uuid.New().String(),
			PartUUIDs:  []string{uuid.New().String()},
			TotalPrice: 50.25,
			Status:     model.OrderStatusPENDINGPAYMENT,
		}

		_, err := repo.CreateOrder(ctx, order)
		require.NoError(t, err)

		// Try to create same order again
		_, err = repo.CreateOrder(ctx, order)
		require.Error(t, err)
	})
}

func TestInMemoryOrderRepository_Get(t *testing.T) {
	repo := NewInMemoryOrderRepository()
	ctx := context.Background()

	t.Run("successfully get order", func(t *testing.T) {
		order := &model.Order{
			OrderUUID:  uuid.New().String(),
			UserUUID:   uuid.New().String(),
			PartUUIDs:  []string{uuid.New().String()},
			TotalPrice: 75.00,
			Status:     model.OrderStatusPENDINGPAYMENT,
		}

		_, err := repo.CreateOrder(ctx, order)
		require.NoError(t, err)

		stored, err := repo.Get(ctx, order.OrderUUID)
		require.NoError(t, err)
		assert.Equal(t, order.OrderUUID, stored.OrderUUID)
	})

	t.Run("fail on non-existent order", func(t *testing.T) {
		_, err := repo.Get(ctx, uuid.New().String())
		require.Error(t, err)
	})
}

func TestInMemoryOrderRepository_UpdateOrder(t *testing.T) {
	repo := NewInMemoryOrderRepository()
	ctx := context.Background()

	t.Run("successfully update order", func(t *testing.T) {
		order := &model.Order{
			OrderUUID:  uuid.New().String(),
			UserUUID:   uuid.New().String(),
			PartUUIDs:  []string{uuid.New().String()},
			TotalPrice: 100.00,
			Status:     model.OrderStatusPENDINGPAYMENT,
		}

		_, err := repo.CreateOrder(ctx, order)
		require.NoError(t, err)

		// Update order
		order.Status = model.OrderStatusPAID
		order.TotalPrice = 150.00

		err = repo.UpdateOrder(ctx, order)
		require.NoError(t, err)

		// Verify updates
		stored, err := repo.Get(ctx, order.OrderUUID)
		require.NoError(t, err)
		assert.Equal(t, repomodel.OrderStatus(order.Status), stored.Status)
		assert.Equal(t, order.TotalPrice, stored.TotalPrice)
	})

	t.Run("fail on non-existent order", func(t *testing.T) {
		order := &model.Order{
			OrderUUID: uuid.New().String(),
			Status:    model.OrderStatusPAID,
		}
		err := repo.UpdateOrder(ctx, order)
		require.Error(t, err)
	})
}

func TestInMemoryOrderRepository_DeleteOrder(t *testing.T) {
	repo := NewInMemoryOrderRepository()
	ctx := context.Background()

	t.Run("successfully delete order", func(t *testing.T) {
		order := &model.Order{
			OrderUUID:  uuid.New().String(),
			UserUUID:   uuid.New().String(),
			PartUUIDs:  []string{uuid.New().String()},
			TotalPrice: 100.00,
			Status:     model.OrderStatusPENDINGPAYMENT,
		}

		_, err := repo.CreateOrder(ctx, order)
		require.NoError(t, err)

		err = repo.DeleteOrder(ctx, order.OrderUUID)
		require.NoError(t, err)

		// Verify order was deleted
		_, err = repo.Get(ctx, order.OrderUUID)
		require.Error(t, err)
	})

	t.Run("fail on non-existent order", func(t *testing.T) {
		err := repo.DeleteOrder(ctx, uuid.New().String())
		require.Error(t, err)
	})
}

func TestInMemoryOrderRepository_CancelOrder(t *testing.T) {
	repo := NewInMemoryOrderRepository()
	ctx := context.Background()

	t.Run("successfully cancel order", func(t *testing.T) {
		order := &model.Order{
			OrderUUID:  uuid.New().String(),
			UserUUID:   uuid.New().String(),
			PartUUIDs:  []string{uuid.New().String()},
			TotalPrice: 100.00,
			Status:     model.OrderStatusPENDINGPAYMENT,
		}

		_, err := repo.CreateOrder(ctx, order)
		require.NoError(t, err)

		err = repo.CancelOrder(ctx, order.OrderUUID)
		require.NoError(t, err)

		// Verify order was cancelled
		stored, err := repo.Get(ctx, order.OrderUUID)
		require.NoError(t, err)
		assert.Equal(t, repomodel.OrderStatusCANCELLED, stored.Status)
	})

	t.Run("fail on non-existent order", func(t *testing.T) {
		err := repo.CancelOrder(ctx, uuid.New().String())
		require.Error(t, err)
	})
}
