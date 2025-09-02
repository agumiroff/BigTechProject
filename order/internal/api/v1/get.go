package api

import (
	"context"
	"log"

	"github.com/agumiroff/BigTechProject/order/v1/internal/converter"
	order_v1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
)

func (a *api) GetOrder(ctx context.Context, req order_v1.GetOrderByUuidParams) (res *order_v1.Order, err error) {
	order, err := a.service.GetOrder(ctx, req.OrderUUID.String())
	if err != nil {
		log.Printf("Failed to create order %v", err)
		return &order_v1.Order{}, err
	}
	return converter.ModelOrderToProto(order), nil
}
