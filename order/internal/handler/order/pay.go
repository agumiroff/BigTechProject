package order

import (
	"context"

	order_v1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
)

func (h *OrderHandler) PayOrder(ctx context.Context, req *order_v1.PayOrderRequest, params order_v1.PayOrderParams) (order_v1.PayOrderRes, error) {
	return h.API.PayOrder(ctx, req, params)
}
