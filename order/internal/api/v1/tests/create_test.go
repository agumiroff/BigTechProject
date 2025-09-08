package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	api "github.com/agumiroff/BigTechProject/order/v1/internal/api/v1"
	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	"github.com/agumiroff/BigTechProject/order/v1/internal/service/mocks"
	order_v1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
)

func TestCreateOrder_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := mocks.NewOrderService(t)
	apiHandler := api.NewAPI(mockService)

	req := &order_v1.CreateOrderRequest{
		UserUUID:  "test-user",
		PartUuids: []string{"part1", "part2"},
	}

	expectedResp := &model.CreateOrderResponse{
		OrderUUID:  "test-uuid",
		TotalPrice: 100.50,
	}

	// Mock service call
	mockService.EXPECT().CreateOrder(ctx, &model.CreateOrderRequest{
		UserUUID:  req.UserUUID,
		PartUUIDs: req.PartUuids,
	}).Return(expectedResp, nil)

	// Act
	resp, err := apiHandler.CreateOrder(ctx, req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, expectedResp.OrderUUID, resp.OrderUUID)
	require.Equal(t, expectedResp.TotalPrice, resp.TotalPrice)
}

func TestCreateOrder_ServiceError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := mocks.NewOrderService(t)
	apiHandler := api.NewAPI(mockService)

	req := &order_v1.CreateOrderRequest{
		UserUUID:  "test-user",
		PartUuids: []string{"part1"},
	}

	expectedErr := errors.New("service error")

	// Mock service error
	mockService.EXPECT().CreateOrder(ctx, &model.CreateOrderRequest{
		UserUUID:  req.UserUUID,
		PartUUIDs: req.PartUuids,
	}).Return(nil, expectedErr)

	// Act
	resp, err := apiHandler.CreateOrder(ctx, req)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, resp)
}

func TestCreateOrder_EmptyRequest(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := mocks.NewOrderService(t)
	apiHandler := api.NewAPI(mockService)

	req := &order_v1.CreateOrderRequest{}

	// Mock service call
	// No mock expectations needed since validation is done in the API layer

	// Act
	resp, err := apiHandler.CreateOrder(ctx, req)

	// Assert
	require.Error(t, err)
	require.Contains(t, err.Error(), "validation error")
	require.Nil(t, resp)
}
