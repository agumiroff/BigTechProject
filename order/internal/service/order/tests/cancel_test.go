package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	mockex "github.com/agumiroff/BigTechProject/order/v1/external/repository/mocks"
	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	"github.com/agumiroff/BigTechProject/order/v1/internal/repository/mocks"
	repomodel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/order/v1/internal/service/order"
	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
)

func TestCancelOrder_Success(t *testing.T) {
	// Arrange
	logger.SetNopLogger()
	ctx := context.Background()
	mockRepo := mocks.NewOrderRepository(t)
	mockExRepo := mockex.NewOrderRepository(t)
	svc := order.NewService(mockRepo, mockExRepo)

	orderUUID := "test-order-uuid"

	existingOrder := &repomodel.OrderRow{
		OrderUUID:  orderUUID,
		UserUUID:   "test-user",
		TotalPrice: 100.0,
		Status:     string(model.OrderStatusPENDINGPAYMENT),
	}
	parts := []string{"part1"}

	// Mock get order
	mockRepo.On("GetOrder", ctx, orderUUID).Return(existingOrder, parts, nil)

	// Mock cancel order
	mockRepo.On("CancelOrder", ctx, orderUUID).Return(nil)

	// No need to mock external publish event as it's not called in the CancelOrder method

	// Act
	err := svc.CancelOrder(ctx, orderUUID)

	// Assert
	require.NoError(t, err)
}

func TestCancelOrder_OrderNotFound(t *testing.T) {
	// Arrange
	logger.SetNopLogger()
	ctx := context.Background()
	mockRepo := mocks.NewOrderRepository(t)
	mockExRepo := mockex.NewOrderRepository(t)
	svc := order.NewService(mockRepo, mockExRepo)

	orderUUID := "test-order-uuid"

	// Mock get order not found
	mockRepo.On("GetOrder", ctx, orderUUID).Return(nil, nil, errors.New("order not found"))

	// Act
	err := svc.CancelOrder(ctx, orderUUID)

	// Assert
	require.Error(t, err)
	require.Contains(t, err.Error(), "order not found")
}

func TestCancelOrder_DeleteError(t *testing.T) {
	// Arrange
	logger.SetNopLogger()
	ctx := context.Background()
	mockRepo := mocks.NewOrderRepository(t)
	mockExRepo := mockex.NewOrderRepository(t)
	svc := order.NewService(mockRepo, mockExRepo)

	orderUUID := "test-order-uuid"
	expectedErr := errors.New("delete error")

	existingOrder := &repomodel.OrderRow{
		OrderUUID:  orderUUID,
		UserUUID:   "test-user",
		TotalPrice: 100.0,
		Status:     string(model.OrderStatusPENDINGPAYMENT),
	}
	parts := []string{"part1"}

	// Mock get order success
	mockRepo.On("GetOrder", ctx, orderUUID).Return(existingOrder, parts, nil)

	// Mock cancel error
	mockRepo.On("CancelOrder", ctx, orderUUID).Return(expectedErr)

	// Act
	err := svc.CancelOrder(ctx, orderUUID)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
}

func TestCancelOrder_OrderAlreadyPaid(t *testing.T) {
	// Arrange
	logger.SetNopLogger()
	ctx := context.Background()
	mockRepo := mocks.NewOrderRepository(t)
	mockExRepo := mockex.NewOrderRepository(t)
	svc := order.NewService(mockRepo, mockExRepo)

	orderUUID := "test-order-uuid"

	existingOrder := &repomodel.OrderRow{
		OrderUUID:  orderUUID,
		UserUUID:   "test-user",
		TotalPrice: 100.0,
		Status:     string(model.OrderStatusPAID),
	}
	parts := []string{"part1"}

	// Mock get order success
	mockRepo.On("GetOrder", ctx, orderUUID).Return(existingOrder, parts, nil)

	// This is a validation error in the service layer, so no mock needed for CancelOrder
	// because the service should return error before calling CancelOrder

	// Act
	err := svc.CancelOrder(ctx, orderUUID)

	// Assert
	require.Error(t, err)
	require.Contains(t, err.Error(), "forbidden")
}
