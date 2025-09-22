package api

import (
	"context"

	orderV1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
)

type OrderAPI interface {
	CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (*orderV1.CreateOrderResponse, error)
	PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (*orderV1.PayOrderResponse, error)
	GetOrder(ctx context.Context, params orderV1.GetOrderByUuidParams) (*orderV1.Order, error)
	CancelOrder(ctx context.Context, params orderV1.CancelOrderByUuidParams) (orderV1.CancelOrderByUuidRes, error)
}
