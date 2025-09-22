package order

import (
	"context"

	order_v1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
)

func (h *OrderHandler) CancelOrderByUuid(ctx context.Context, params order_v1.CancelOrderByUuidParams) (order_v1.CancelOrderByUuidRes, error) {
	res, err := h.API.CancelOrder(ctx, params)
	if err != nil {
		return nil, err
	}

	return res, nil
}
