package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	api "github.com/agumiroff/BigTechProject/order/v1/internal/api/v1"
	rModel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/order/v1/internal/service/mocks"
	order_v1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
)

func TestCancelOrder_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := mocks.NewOrderService(t)
	apiHandler := api.NewAPI(mockService)

	orderID := uuid.New()
	params := order_v1.CancelOrderByUuidParams{
		OrderUUID: orderID,
	}

	// Mock service call
	mockService.EXPECT().CancelOrder(ctx, orderID.String()).Return(nil)

	// Act
	resp, err := apiHandler.CancelOrder(ctx, params)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	response := resp.(*order_v1.CancelOrderResponse)
	require.Equal(t, orderID, response.OrderUUID)
}

func TestCancelOrder_OrderNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := mocks.NewOrderService(t)
	apiHandler := api.NewAPI(mockService)

	orderID := uuid.New()
	params := order_v1.CancelOrderByUuidParams{
		OrderUUID: orderID,
	}

	// Mock service error
	mockService.EXPECT().CancelOrder(ctx, orderID.String()).Return(rModel.ErrOrderNotFound)

	// Act
	resp, err := apiHandler.CancelOrder(ctx, params)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, rModel.ErrOrderNotFound)
	require.Nil(t, resp)
}

func TestCancelOrder_ServiceError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := mocks.NewOrderService(t)
	apiHandler := api.NewAPI(mockService)

	orderID := uuid.New()
	params := order_v1.CancelOrderByUuidParams{
		OrderUUID: orderID,
	}

	expectedErr := errors.New("service error")

	// Mock service error
	mockService.EXPECT().CancelOrder(ctx, orderID.String()).Return(expectedErr)

	// Act
	resp, err := apiHandler.CancelOrder(ctx, params)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, resp)
}
