package repository

import (
	"context"

	repomodel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *repomodel.OrderRow, parts []string) (string, error)
	GetOrder(ctx context.Context, uuid string) (order *repomodel.OrderRow, parts []string, err error)
	UpdateOrder(ctx context.Context, order *repomodel.OrderRow) error
	DeleteOrder(ctx context.Context, uuid string) error
	CancelOrder(ctx context.Context, uuid string) error
}
