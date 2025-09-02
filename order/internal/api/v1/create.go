package api

import (
	"context"
	"log"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	order_v1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (res *order_v1.CreateOrderResponse, err error) {
	order, err := a.service.CreateOrder(ctx, &model.CreateOrderRequest{
		UserUUID:  req.UserUUID,
		PartUUIDs: req.PartUuids,
	})
	if err != nil {
		log.Printf("Failed to create order %v", err)
		return &order_v1.CreateOrderResponse{}, err
	}

	return &order_v1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}
