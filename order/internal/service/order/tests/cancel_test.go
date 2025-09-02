package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	mockex "github.com/agumiroff/BigTechProject/order/v1/external/repository/mocks"
	"github.com/agumiroff/BigTechProject/order/v1/internal/repository/mocks"
	rModel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/order/v1/internal/service/order"
)

func TestCancelOrder_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := mocks.NewOrderRepository(t)
	mockExRepo := mockex.NewOrderRepository(t)
	svc := order.NewService(mockRepo, mockExRepo)

	orderUUID := "test-order-uuid"

	existingOrder := &rModel.Order{
		OrderUUID:  orderUUID,
		UserUUID:   "test-user",
		PartUUIDs:  []string{"part1"},
		TotalPrice: 100.0,
		Status:     rModel.OrderStatusPENDINGPAYMENT,
	}

	// Mock get order
	mockRepo.On("Get", orderUUID).Return(existingOrder, nil)

	// Mock delete order
	mockRepo.On("DeleteOrder", orderUUID).Return(nil)

	// Act
	err := svc.CancelOrder(ctx, orderUUID)

	// Assert
	require.NoError(t, err)
}

func TestCancelOrder_OrderNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := mocks.NewOrderRepository(t)
	mockExRepo := mockex.NewOrderRepository(t)
	svc := order.NewService(mockRepo, mockExRepo)

	orderUUID := "test-order-uuid"

	// Mock get order not found
	mockRepo.On("Get", orderUUID).Return(nil, rModel.ErrOrderNotFound)

	// Act
	err := svc.CancelOrder(ctx, orderUUID)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, rModel.ErrOrderNotFound)
}

func TestCancelOrder_DeleteError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := mocks.NewOrderRepository(t)
	mockExRepo := mockex.NewOrderRepository(t)
	svc := order.NewService(mockRepo, mockExRepo)

	orderUUID := "test-order-uuid"
	expectedErr := errors.New("delete error")

	existingOrder := &rModel.Order{
		OrderUUID:  orderUUID,
		UserUUID:   "test-user",
		PartUUIDs:  []string{"part1"},
		TotalPrice: 100.0,
		Status:     rModel.OrderStatusPENDINGPAYMENT,
	}

	// Mock get order success
	mockRepo.On("Get", orderUUID).Return(existingOrder, nil)

	// Mock delete error
	mockRepo.On("DeleteOrder", orderUUID).Return(expectedErr)

	// Act
	err := svc.CancelOrder(ctx, orderUUID)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
}

func TestCancelOrder_OrderAlreadyPaid(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := mocks.NewOrderRepository(t)
	mockExRepo := mockex.NewOrderRepository(t)
	svc := order.NewService(mockRepo, mockExRepo)

	orderUUID := "test-order-uuid"

	existingOrder := &rModel.Order{
		OrderUUID:  orderUUID,
		UserUUID:   "test-user",
		PartUUIDs:  []string{"part1"},
		TotalPrice: 100.0,
		Status:     rModel.OrderStatusPAID,
	}

	// Mock get order success
	mockRepo.On("Get", orderUUID).Return(existingOrder, nil)

	// Mock delete error
	mockRepo.On("DeleteOrder", orderUUID).Return(rModel.ErrOrderAlreadyPaid)

	// Act
	err := svc.CancelOrder(ctx, orderUUID)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, rModel.ErrOrderAlreadyPaid)
}
