package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	api "github.com/agumiroff/BigTechProject/order/v1/internal/api/v1"
	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	"github.com/agumiroff/BigTechProject/order/v1/internal/service/mocks"
	order_v1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
)

func TestGetOrder_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := mocks.NewOrderService(t)
	apiHandler := api.NewAPI(mockService)

	orderID := uuid.New()
	txID := uuid.New()
	params := order_v1.GetOrderByUuidParams{
		OrderUUID: orderID,
	}

	expectedOrder := &model.Order{
		OrderUUID:       orderID.String(),
		UserUUID:        "test-user",
		PartUUIDs:       []string{"part1", "part2"},
		TotalPrice:      100.50,
		TransactionUUID: txID.String(),
		PaymentMethod:   model.PaymentMethodCARD,
		Status:          model.OrderStatusPAID,
	}

	// Mock service call
	mockService.EXPECT().GetOrder(ctx, orderID.String()).Return(expectedOrder, nil)

	// Act
	resp, err := apiHandler.GetOrder(ctx, params)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, expectedOrder.OrderUUID, resp.OrderUUID)
	require.Equal(t, expectedOrder.UserUUID, resp.UserUUID)
	require.Equal(t, expectedOrder.PartUUIDs, resp.PartUuids)
	require.Equal(t, expectedOrder.TotalPrice, resp.TotalPrice)
	require.Equal(t, order_v1.OrderStatus(expectedOrder.Status), resp.Status)

	txUUID, ok := resp.TransactionUUID.Get()
	require.True(t, ok)
	require.Equal(t, expectedOrder.TransactionUUID, txUUID)

	payMethod, ok := resp.PaymentMethod.Get()
	require.True(t, ok)
	require.Equal(t, order_v1.PaymentMethod(expectedOrder.PaymentMethod), payMethod)
}

func TestGetOrder_ServiceError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := mocks.NewOrderService(t)
	apiHandler := api.NewAPI(mockService)

	orderID := uuid.New()
	params := order_v1.GetOrderByUuidParams{
		OrderUUID: orderID,
	}

	expectedErr := errors.New("service error")

	// Mock service error
	mockService.EXPECT().GetOrder(ctx, orderID.String()).Return(nil, expectedErr)

	// Act
	resp, err := apiHandler.GetOrder(ctx, params)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, resp)
}

func TestGetOrder_InvalidTransactionUUID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := mocks.NewOrderService(t)
	apiHandler := api.NewAPI(mockService)

	orderID := uuid.New()
	params := order_v1.GetOrderByUuidParams{
		OrderUUID: orderID,
	}

	expectedOrder := &model.Order{
		OrderUUID:       orderID.String(),
		UserUUID:        "test-user",
		TransactionUUID: "invalid-uuid",
	}

	// Mock service call
	mockService.EXPECT().GetOrder(ctx, orderID.String()).Return(expectedOrder, nil)

	// Act
	resp, err := apiHandler.GetOrder(ctx, params)

	// Assert
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid UUID")
	require.Nil(t, resp)
}
