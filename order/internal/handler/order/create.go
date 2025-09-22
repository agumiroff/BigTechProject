package order

import (
	"context"

	order_v1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
)

func (h *OrderHandler) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {
	return h.API.CreateOrder(ctx, req)
}
