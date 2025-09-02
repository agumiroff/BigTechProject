package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mockex "github.com/agumiroff/BigTechProject/order/v1/external/repository/mocks"
	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	"github.com/agumiroff/BigTechProject/order/v1/internal/repository/mocks"
	rModel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/order/v1/internal/service/order"
	paymentv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

func TestPayOrder_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := mocks.NewOrderRepository(t)
	mockExRepo := mockex.NewOrderRepository(t)
	svc := order.NewService(mockRepo, mockExRepo)

	orderUUID := "test-order-uuid"
	transactionUUID := "test-transaction-uuid"

	req := &model.PayOrderRequest{
		OrderUUID:     orderUUID,
		PaymentMethod: model.PaymentMethodCARD,
	}

	existingOrder := &rModel.Order{
		OrderUUID:  orderUUID,
		UserUUID:   "test-user",
		PartUUIDs:  []string{"part1"},
		TotalPrice: 100.0,
		Status:     rModel.OrderStatusPENDINGPAYMENT,
	}

	// Mock payment service response
	mockExRepo.On("PayOrder", ctx, &paymentv1.PayOrderRequest{
		Payment: &paymentv1.Payment{
			OrderUuid:     orderUUID,
			PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
		},
	}).Return(&paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil)

	// Mock get order
	mockRepo.On("Get", orderUUID).Return(existingOrder, nil)

	// Mock update order
	mockRepo.On("UpdateOrder", ctx, mock.MatchedBy(func(order *model.Order) bool {
		if order == nil {
			return false
		}
		return order.OrderUUID == orderUUID &&
			order.Status == model.OrderStatusPAID &&
			order.PaymentMethod == model.PaymentMethodCARD &&
			order.TransactionUUID == transactionUUID &&
			order.UserUUID == existingOrder.UserUUID &&
			order.TotalPrice == existingOrder.TotalPrice &&
			len(order.PartUUIDs) == len(existingOrder.PartUUIDs) &&
			order.PartUUIDs[0] == existingOrder.PartUUIDs[0]
	})).Return(nil)

	// Act
	resp, err := svc.PayOrder(ctx, req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, transactionUUID, resp.TransactionUUID)

	mockRepo.AssertExpectations(t)
	mockExRepo.AssertExpectations(t)
}

func TestPayOrder_PaymentError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := mocks.NewOrderRepository(t)
	mockExRepo := mockex.NewOrderRepository(t)
	svc := order.NewService(mockRepo, mockExRepo)

	orderUUID := "test-order-uuid"
	expectedErr := errors.New("payment service error")

	req := &model.PayOrderRequest{
		OrderUUID:     orderUUID,
		PaymentMethod: model.PaymentMethodCARD,
	}

	// Mock payment service error
	mockExRepo.On("PayOrder", ctx, mock.Anything).Return(nil, expectedErr)

	// Act
	resp, err := svc.PayOrder(ctx, req)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Equal(t, &model.PayOrderResponse{}, resp)

	mockRepo.AssertExpectations(t)
	mockExRepo.AssertExpectations(t)
}

func TestPayOrder_OrderNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := mocks.NewOrderRepository(t)
	mockExRepo := mockex.NewOrderRepository(t)
	svc := order.NewService(mockRepo, mockExRepo)

	orderUUID := "test-order-uuid"
	transactionUUID := "test-transaction-uuid"

	req := &model.PayOrderRequest{
		OrderUUID:     orderUUID,
		PaymentMethod: model.PaymentMethodCARD,
	}

	// Mock payment service success
	mockExRepo.On("PayOrder", ctx, mock.Anything).Return(&paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil)

	// Mock order not found
	mockRepo.On("Get", orderUUID).Return(nil, rModel.ErrOrderNotFound)

	// Act
	resp, err := svc.PayOrder(ctx, req)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, rModel.ErrOrderNotFound)
	require.Equal(t, &model.PayOrderResponse{}, resp)

	mockRepo.AssertExpectations(t)
	mockExRepo.AssertExpectations(t)
}

func TestPayOrder_UpdateError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := mocks.NewOrderRepository(t)
	mockExRepo := mockex.NewOrderRepository(t)
	svc := order.NewService(mockRepo, mockExRepo)

	orderUUID := "test-order-uuid"
	transactionUUID := "test-transaction-uuid"

	req := &model.PayOrderRequest{
		OrderUUID:     orderUUID,
		PaymentMethod: model.PaymentMethodCARD,
	}

	existingOrder := &rModel.Order{
		OrderUUID: orderUUID,
		Status:    rModel.OrderStatusPENDINGPAYMENT,
	}

	expectedErr := errors.New("update error")

	// Mock payment service success
	mockExRepo.On("PayOrder", ctx, mock.Anything).Return(&paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil)

	// Mock get order success
	mockRepo.On("Get", orderUUID).Return(existingOrder, nil)

	// Mock update error
	mockRepo.On("UpdateOrder", ctx, mock.Anything).Return(expectedErr)

	// Act
	resp, err := svc.PayOrder(ctx, req)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Equal(t, &model.PayOrderResponse{}, resp)

	mockRepo.AssertExpectations(t)
	mockExRepo.AssertExpectations(t)
}
