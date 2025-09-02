package handler

import (
	"context"

	order_v1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
)

type OrderHandler interface {
	CancelOrderByUuid(ctx context.Context, params order_v1.CancelOrderByUuidParams) (order_v1.CancelOrderByUuidRes, error)
	CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error)
	GetOrderByUuid(ctx context.Context, params order_v1.GetOrderByUuidParams) (order_v1.Order, error)
	PayOrder(ctx context.Context, req *order_v1.PayOrderRequest, params order_v1.PayOrderParams) (order_v1.PayOrderRes, error)
	NewError(ctx context.Context, err error) *order_v1.GenericErrorStatusCode
}
