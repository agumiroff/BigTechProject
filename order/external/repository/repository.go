package repository

import (
	"context"

	invV1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
	payV1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

type OrderRepository interface {
	ListParts(ctx context.Context, req *invV1.ListPartsRequest) (*invV1.ListPartsResponse, error)
	PayOrder(ctx context.Context, req *payV1.PayOrderRequest) (*payV1.PayOrderResponse, error)
}
