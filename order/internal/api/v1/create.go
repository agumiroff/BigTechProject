package api

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
	order_v1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (res *order_v1.CreateOrderResponse, err error) {
	if req == nil || req.UserUUID == "" || len(req.PartUuids) == 0 {
		return nil, errors.New("validation error: missing required fields")
	}

	order, err := a.service.CreateOrder(ctx, &model.CreateOrderRequest{
		UserUUID:  req.UserUUID,
		PartUUIDs: req.PartUuids,
	})
	if err != nil {
		logger.Error(ctx, "Failed to create order", zap.Error(err))
		return nil, err
	}

	return &order_v1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}
