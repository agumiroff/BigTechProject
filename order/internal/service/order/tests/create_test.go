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
	repomodel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/order/v1/internal/service/order"
	inventoryv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
)

func TestCreateOrder_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := mocks.NewOrderRepository(t)
	mockExRepo := mockex.NewOrderRepository(t)
	svc := order.NewService(mockRepo, mockExRepo)

	req := &model.CreateOrderRequest{
		UserUUID:  "test-user",
		PartUUIDs: []string{"part1", "part2"},
	}

	parts := []*inventoryv1.Part{
		{
			Uuid:  "part1",
			Price: 50.0,
		},
		{
			Uuid:  "part2",
			Price: 75.0,
		},
	}

	// Mock external repo response
	mockExRepo.On("ListParts", ctx, &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartsFilter{
			Uuids: req.PartUUIDs,
		},
	}).Return(&inventoryv1.ListPartsResponse{
		Parts: parts,
	}, nil)

	// Mock internal repo response
	mockRepo.On("CreateOrder", mock.Anything, mock.MatchedBy(func(orderRow *repomodel.OrderRow) bool {
		return orderRow.UserUUID == req.UserUUID &&
			orderRow.TotalPrice == 125.0 &&
			orderRow.Status == string(model.OrderStatusPENDINGPAYMENT)
	}), mock.MatchedBy(func(parts []string) bool {
		return len(parts) == len(req.PartUUIDs)
	})).Return("test-uuid", nil)

	// Act
	resp, err := svc.CreateOrder(ctx, req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 125.0, resp.TotalPrice)
	require.NotEmpty(t, resp.OrderUUID)
}

func TestCreateOrder_PartNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := mocks.NewOrderRepository(t)
	mockExRepo := mockex.NewOrderRepository(t)
	svc := order.NewService(mockRepo, mockExRepo)

	req := &model.CreateOrderRequest{
		UserUUID:  "test-user",
		PartUUIDs: []string{"part1", "part2"},
	}

	// Only one part exists
	parts := []*inventoryv1.Part{
		{
			Uuid:  "part1",
			Price: 50.0,
		},
	}

	// Mock external repo response
	mockExRepo.On("ListParts", ctx, &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartsFilter{
			Uuids: req.PartUUIDs,
		},
	}).Return(&inventoryv1.ListPartsResponse{
		Parts: parts,
	}, nil)

	// Act
	resp, err := svc.CreateOrder(ctx, req)

	// Assert
	require.Error(t, err)
	require.Nil(t, resp)
	require.Contains(t, err.Error(), "parts not found")
}

func TestCreateOrder_ListPartsError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := mocks.NewOrderRepository(t)
	mockExRepo := mockex.NewOrderRepository(t)
	svc := order.NewService(mockRepo, mockExRepo)

	req := &model.CreateOrderRequest{
		UserUUID:  "test-user",
		PartUUIDs: []string{"part1"},
	}

	expectedErr := errors.New("inventory service error")

	// Mock external repo response
	mockExRepo.On("ListParts", ctx, &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartsFilter{
			Uuids: req.PartUUIDs,
		},
	}).Return(nil, expectedErr)

	// Act
	resp, err := svc.CreateOrder(ctx, req)

	// Assert
	require.Error(t, err)
	require.Nil(t, resp)
	require.ErrorIs(t, err, expectedErr)
}
