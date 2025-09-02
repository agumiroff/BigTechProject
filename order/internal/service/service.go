package service

import (
	"context"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
)

type OrderService interface {
	PayOrder(ctx context.Context, req *model.PayOrderRequest) (res *model.PayOrderResponse, err error)
	CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (res *model.CreateOrderResponse, err error)
	GetOrder(ctx context.Context, uuid string) (*model.Order, error)
	CancelOrder(ctx context.Context, uuid string) error
	DeleteOrder(ctx context.Context, uuid string) error
}
