package order

import (
	"context"

	order_v1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
)

func (h *OrderHandler) GetOrderByUuid(ctx context.Context, params order_v1.GetOrderByUuidParams) (order_v1.GetOrderByUuidRes, error) {
	return h.API.GetOrder(ctx, params)
}
