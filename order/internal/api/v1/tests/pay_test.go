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

func TestPayOrder_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := mocks.NewOrderService(t)
	apiHandler := api.NewAPI(mockService)

	orderID := uuid.New()
	txID := uuid.New()
	params := order_v1.PayOrderParams{
		OrderUUID: orderID,
	}

	req := &order_v1.PayOrderRequest{
		PaymentMethod: order_v1.PaymentMethodCARD,
	}

	expectedResp := &model.PayOrderResponse{
		TransactionUUID: txID,
	}

	// Mock service call
	mockService.EXPECT().PayOrder(ctx, &model.PayOrderRequest{
		OrderUUID:     orderID.String(),
		PaymentMethod: model.PaymentMethodCARD,
	}).Return(expectedResp, nil)

	// Act
	resp, err := apiHandler.PayOrder(ctx, req, params)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, txID, resp.TransactionUUID)
}

func TestPayOrder_ServiceError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := mocks.NewOrderService(t)
	apiHandler := api.NewAPI(mockService)

	orderID := uuid.New()
	params := order_v1.PayOrderParams{
		OrderUUID: orderID,
	}

	req := &order_v1.PayOrderRequest{
		PaymentMethod: order_v1.PaymentMethodCARD,
	}

	expectedErr := errors.New("service error")

	// Mock service error
	mockService.EXPECT().PayOrder(ctx, &model.PayOrderRequest{
		OrderUUID:     orderID.String(),
		PaymentMethod: model.PaymentMethodCARD,
	}).Return(nil, expectedErr)

	// Act
	resp, err := apiHandler.PayOrder(ctx, req, params)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, resp)
}

func TestPayOrder_EmptyRequest(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := mocks.NewOrderService(t)
	apiHandler := api.NewAPI(mockService)

	orderID := uuid.New()
	params := order_v1.PayOrderParams{
		OrderUUID: orderID,
	}

	req := &order_v1.PayOrderRequest{}

	// No mock expectations needed since validation is done in the API layer

	// Act
	resp, err := apiHandler.PayOrder(ctx, req, params)

	// Assert
	require.Error(t, err)
	require.Contains(t, err.Error(), "validation error")
	require.Nil(t, resp)
}

func TestPayOrder_InvalidTransactionUUID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockService := mocks.NewOrderService(t)
	apiHandler := api.NewAPI(mockService)

	orderID := uuid.New()
	params := order_v1.PayOrderParams{
		OrderUUID: orderID,
	}

	req := &order_v1.PayOrderRequest{
		PaymentMethod: order_v1.PaymentMethodCARD,
	}

	// Return an error during UUID parsing
	mockService.EXPECT().PayOrder(ctx, &model.PayOrderRequest{
		OrderUUID:     orderID.String(),
		PaymentMethod: model.PaymentMethodCARD,
	}).Return(nil, errors.New("invalid UUID format"))

	// Act
	resp, err := apiHandler.PayOrder(ctx, req, params)

	// Assert
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid UUID")
	require.Nil(t, resp)
}
