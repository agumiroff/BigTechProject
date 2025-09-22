package order

import (
	"github.com/agumiroff/BigTechProject/order/v1/internal/api"
	order_v1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
)

type OrderHandler struct {
	API api.OrderAPI
}

var _ order_v1.Handler = (*OrderHandler)(nil)

func NewHandler(api api.OrderAPI) *OrderHandler {
	return &OrderHandler{
		API: api,
	}
}
