package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mockex "github.com/agumiroff/BigTechProject/order/v1/external/repository/mocks"
	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	"github.com/agumiroff/BigTechProject/order/v1/internal/repository/mocks"
	repomodel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/order/v1/internal/service/order"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
	paymentv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

func TestPayOrder_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := mocks.NewOrderRepository(t)
	mockExRepo := mockex.NewOrderRepository(t)
	svc := order.NewService(mockRepo, mockExRepo)

	orderUUID := "550e8400-e29b-41d4-a716-446655440001"
	transactionUUID := "550e8400-e29b-41d4-a716-446655440002"

	req := &model.PayOrderRequest{
		OrderUUID:     orderUUID,
		PaymentMethod: model.PaymentMethodCARD,
	}

	existingOrder := &repomodel.Order{
		OrderUUID:  orderUUID,
		UserUUID:   "test-user",
		PartUUIDs:  []string{"part1"},
		TotalPrice: 100.0,
		Status:     repomodel.OrderStatusPENDINGPAYMENT,
	}

	// Mock get order (initial status check)
	mockRepo.On("Get", ctx, orderUUID).Return(existingOrder, nil)

	// Mock payment service response
	mockExRepo.On("PayOrder", ctx, mock.MatchedBy(func(req *paymentv1.PayOrderRequest) bool {
		if req == nil || req.Payment == nil {
			return false
		}
		return req.Payment.OrderUuid == orderUUID &&
			req.Payment.UserUuid == existingOrder.UserUUID &&
			req.Payment.PaymentMethod == paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	})).Return(&paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil)

	// Mock update order
	mockRepo.On("UpdateOrder", ctx, mock.MatchedBy(func(order *model.Order) bool {
		return order != nil &&
			order.OrderUUID == orderUUID &&
			order.Status == model.OrderStatusPAID &&
			order.PaymentMethod == model.PaymentMethodCARD &&
			order.TransactionUUID == transactionUUID &&
			order.UserUUID == existingOrder.UserUUID &&
			order.TotalPrice == existingOrder.TotalPrice &&
			len(order.PartUUIDs) == len(existingOrder.PartUUIDs) &&
			order.PartUUIDs[0] == existingOrder.PartUUIDs[0]
	})).Return(nil)

	// Mock publish event
	mockExRepo.On("PublishOrderEvent", ctx, orderUUID, existingOrder.UserUUID).Return(nil)

	// Act
	resp, err := svc.PayOrder(ctx, req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)

	expectedUUID, err := uuid.Parse(transactionUUID)
	require.NoError(t, err)
	require.Equal(t, expectedUUID, resp.TransactionUUID)

	mockRepo.AssertExpectations(t)
	mockExRepo.AssertExpectations(t)
}

func TestPayOrder_PaymentError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := mocks.NewOrderRepository(t)
	mockExRepo := mockex.NewOrderRepository(t)
	svc := order.NewService(mockRepo, mockExRepo)

	orderUUID := "550e8400-e29b-41d4-a716-446655440001"

	req := &model.PayOrderRequest{
		OrderUUID:     orderUUID,
		PaymentMethod: model.PaymentMethodCARD,
	}

	// Mock get order not found
	mockRepo.On("Get", ctx, orderUUID).Return(nil, apperrors.ErrNotFound)

	// Act
	resp, err := svc.PayOrder(ctx, req)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, apperrors.ErrNotFound)
	require.Nil(t, resp)

	mockRepo.AssertExpectations(t)
	mockExRepo.AssertExpectations(t)
}

func TestPayOrder_OrderNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := mocks.NewOrderRepository(t)
	mockExRepo := mockex.NewOrderRepository(t)
	svc := order.NewService(mockRepo, mockExRepo)

	orderUUID := "550e8400-e29b-41d4-a716-446655440001"

	req := &model.PayOrderRequest{
		OrderUUID:     orderUUID,
		PaymentMethod: model.PaymentMethodCARD,
	}

	// Mock get order not found
	mockRepo.On("Get", ctx, orderUUID).Return(nil, apperrors.ErrNotFound)

	// Act
	resp, err := svc.PayOrder(ctx, req)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, apperrors.ErrNotFound)
	require.Nil(t, resp)

	mockRepo.AssertExpectations(t)
	mockExRepo.AssertExpectations(t)
}

func TestPayOrder_UpdateError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := mocks.NewOrderRepository(t)
	mockExRepo := mockex.NewOrderRepository(t)
	svc := order.NewService(mockRepo, mockExRepo)

	orderUUID := "550e8400-e29b-41d4-a716-446655440001"
	transactionUUID := "550e8400-e29b-41d4-a716-446655440002"

	req := &model.PayOrderRequest{
		OrderUUID:     orderUUID,
		PaymentMethod: model.PaymentMethodCARD,
	}

	existingOrder := &repomodel.Order{
		OrderUUID:  orderUUID,
		UserUUID:   "test-user",
		PartUUIDs:  []string{"part1"},
		TotalPrice: 100.0,
		Status:     repomodel.OrderStatusPENDINGPAYMENT,
	}

	expectedErr := errors.New("update error")

	// Mock get order success
	mockRepo.On("Get", ctx, orderUUID).Return(existingOrder, nil)

	// Mock payment service response
	mockExRepo.On("PayOrder", ctx, mock.MatchedBy(func(req *paymentv1.PayOrderRequest) bool {
		if req == nil || req.Payment == nil {
			return false
		}
		return req.Payment.OrderUuid == orderUUID &&
			req.Payment.UserUuid == existingOrder.UserUUID &&
			req.Payment.PaymentMethod == paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	})).Return(&paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil)

	// Mock update error
	mockRepo.On("UpdateOrder", ctx, mock.MatchedBy(func(order *model.Order) bool {
		return order != nil &&
			order.OrderUUID == orderUUID &&
			order.Status == model.OrderStatusPAID &&
			order.PaymentMethod == model.PaymentMethodCARD &&
			order.TransactionUUID == transactionUUID &&
			order.UserUUID == existingOrder.UserUUID &&
			order.TotalPrice == existingOrder.TotalPrice &&
			len(order.PartUUIDs) == len(existingOrder.PartUUIDs) &&
			order.PartUUIDs[0] == existingOrder.PartUUIDs[0]
	})).Return(expectedErr)

	// Act
	resp, err := svc.PayOrder(ctx, req)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, resp)

	mockRepo.AssertExpectations(t)
	mockExRepo.AssertExpectations(t)
}
