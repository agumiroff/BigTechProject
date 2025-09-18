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

func newTestOrder() *model.Order {
	return &model.Order{
		OrderUUID:  "test-order-uuid",
		UserUUID:   "test-user-uuid",
		PartUUIDs:  []string{"part1", "part2"},
		TotalPrice: 100.50,
		Status:     model.OrderStatusPENDINGPAYMENT,
	}
}

func TestCreateOrderInMemory_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := inmemory.NewInMemoryOrderRepository()
	testOrder := newTestOrder()

	// Act
	resp, err := repo.CreateOrder(ctx, testOrder)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, testOrder.OrderUUID, resp.OrderUUID)
	assert.Equal(t, testOrder.TotalPrice, resp.TotalPrice)

	// Verify order was stored
	stored, err := repo.Get(ctx, testOrder.OrderUUID)
	require.NoError(t, err)
	require.NotNil(t, stored)
	assert.Equal(t, testOrder.UserUUID, stored.UserUUID)
	assert.Equal(t, testOrder.PartUUIDs, stored.PartUUIDs)
	assert.Equal(t, string(testOrder.Status), string(stored.Status))
}

func TestCreateOrderInMemory_DuplicateUUID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := inmemory.NewInMemoryOrderRepository()
	existingOrder := &model.Order{
		OrderUUID:  "same-uuid",
		UserUUID:   "old-user",
		PartUUIDs:  []string{"old-part"},
		TotalPrice: 50.00,
		Status:     model.OrderStatusPENDINGPAYMENT,
	}

	// Create initial order
	_, err := repo.CreateOrder(ctx, existingOrder)
	require.NoError(t, err)

	// Try to create another order with same UUID
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
	require.Error(t, err)
	assert.Nil(t, resp)

	// Verify original order is unchanged
	stored, err := repo.Get(ctx, existingOrder.OrderUUID)
	require.NoError(t, err)
	require.NotNil(t, stored)
	assert.Equal(t, existingOrder.UserUUID, stored.UserUUID)
	assert.Equal(t, existingOrder.PartUUIDs, stored.PartUUIDs)
	assert.Equal(t, string(existingOrder.Status), string(stored.Status))
}

func TestCreateOrderInMemory_NilOrder(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := inmemory.NewInMemoryOrderRepository()

	// Act
	resp, err := repo.CreateOrder(ctx, nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestCreateOrderInMemory_EmptyOrderUUID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := inmemory.NewInMemoryOrderRepository()
	invalidOrder := &model.Order{
		UserUUID:   "test-user",
		PartUUIDs:  []string{"part1"},
		TotalPrice: 100,
		Status:     model.OrderStatusPENDINGPAYMENT,
	}

	// Act
	resp, err := repo.CreateOrder(ctx, invalidOrder)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestCreateOrderInMemory_EmptyUserUUID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := inmemory.NewInMemoryOrderRepository()
	invalidOrder := &model.Order{
		OrderUUID:  "test-uuid",
		PartUUIDs:  []string{"part1"},
		TotalPrice: 100,
		Status:     model.OrderStatusPENDINGPAYMENT,
	}

	// Act
	resp, err := repo.CreateOrder(ctx, invalidOrder)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestCreateOrderInMemory_EmptyPartUUIDs(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := inmemory.NewInMemoryOrderRepository()
	invalidOrder := &model.Order{
		OrderUUID:  "test-uuid",
		UserUUID:   "test-user",
		TotalPrice: 100,
		Status:     model.OrderStatusPENDINGPAYMENT,
	}

	// Act
	resp, err := repo.CreateOrder(ctx, invalidOrder)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestCreateOrderInMemory_ValidationErrors(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := inmemory.NewInMemoryOrderRepository()

	tests := []struct {
		name        string
		order       *model.Order
		expectedErr error
	}{
		{
			name:        "nil order",
			order:       nil,
			expectedErr: apperrors.ErrInvalidRequest,
		},
		{
			name: "empty order UUID",
			order: &model.Order{
				UserUUID:   "test-user",
				PartUUIDs:  []string{"part1"},
				TotalPrice: 100,
				Status:     model.OrderStatusPENDINGPAYMENT,
			},
			expectedErr: apperrors.ErrInvalidRequest,
		},
		{
			name: "empty user UUID",
			order: &model.Order{
				OrderUUID:  "test-uuid",
				PartUUIDs:  []string{"part1"},
				TotalPrice: 100,
				Status:     model.OrderStatusPENDINGPAYMENT,
			},
			expectedErr: apperrors.ErrInvalidRequest,
		},
		{
			name: "empty part UUIDs",
			order: &model.Order{
				OrderUUID:  "test-uuid",
				UserUUID:   "test-user",
				TotalPrice: 100,
				Status:     model.OrderStatusPENDINGPAYMENT,
			},
			expectedErr: apperrors.ErrInvalidRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			resp, err := repo.CreateOrder(ctx, tt.order)

			// Assert
			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.expectedErr)
			assert.Nil(t, resp)
		})
	}
}
