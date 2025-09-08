package repository

import (
	"context"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	repomodel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *model.Order) (*model.CreateOrderResponse, error)
	Get(ctx context.Context, uuid string) (*repomodel.Order, error)
	UpdateOrder(ctx context.Context, order *model.Order) error
	DeleteOrder(ctx context.Context, uuid string) error
	CancelOrder(ctx context.Context, uuid string) error
}
