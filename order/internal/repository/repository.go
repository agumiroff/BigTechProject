package repository

import (
	"context"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	rModel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
)

type OrderRepository interface {
	CreateOrder(req *model.Order) (res *model.CreateOrderResponse)
	Get(uuid string) (*rModel.Order, error)
	UpdateOrder(ctx context.Context, m *model.Order) error
	DeleteOrder(uuid string) error
}
