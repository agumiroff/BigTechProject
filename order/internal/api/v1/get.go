package api

import (
	"context"
	"log"

	"github.com/agumiroff/BigTechProject/order/v1/internal/converter"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
	order_v1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
)

func (a *api) GetOrder(ctx context.Context, req order_v1.GetOrderByUuidParams) (*order_v1.Order, error) {
	order, err := a.service.GetOrder(ctx, req.OrderUUID.String())
	if err != nil {
		log.Printf("Failed to get order %v", err)
		return nil, apperrors.Map(err)
	}

	return converter.ToProtoOrder(order), nil
}
