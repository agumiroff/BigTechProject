package api

import (
	"context"

	"go.uber.org/zap"

	"github.com/agumiroff/BigTechProject/order/v1/internal/converter"
	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
	order_v1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
)

func (a *api) GetOrder(ctx context.Context, req order_v1.GetOrderByUuidParams) (*order_v1.Order, error) {
	order, err := a.service.GetOrder(ctx, req.OrderUUID.String())
	if err != nil {
		logger.Error(ctx, "Failed to get order", zap.Error(err), zap.String("order_uuid", req.OrderUUID.String()))
		return nil, err
	}

	resp, err := converter.ToProtoOrder(order)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
