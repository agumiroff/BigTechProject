package api

import (
	"context"

	order_v1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
)

func (a *api) CancelOrder(ctx context.Context, params order_v1.CancelOrderByUuidParams) (order_v1.CancelOrderByUuidRes, error) {
	err := a.service.CancelOrder(ctx, params.OrderUUID.String())
	if err != nil {
		return &order_v1.CancelOrderResponse{}, err
	}

	return &order_v1.CancelOrderResponse{
		OrderUUID: params.OrderUUID,
	}, nil
}
