package order

import (
	"context"

	order_v1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
)

func (h *OrderHandler) NewError(ctx context.Context, err error) *order_v1.GenericErrorStatusCode {
	return nil
}
