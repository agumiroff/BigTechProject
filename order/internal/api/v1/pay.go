package api

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	orderV1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (*orderV1.PayOrderResponse, error) {
	r, err := a.service.PayOrder(ctx, &model.PayOrderRequest{
		OrderUUID:     params.OrderUUID.String(),
		PaymentMethod: model.PaymentMethod(req.GetPaymentMethod()),
	})
	if err != nil {
		log.Printf("Payment failed, %v", err)
		return &orderV1.PayOrderResponse{}, err
	}

	uuid, err := uuid.Parse(r.TransactionUUID)
	if err != nil {
		log.Printf("failed to parse uuid: %v", err)
		return &orderV1.PayOrderResponse{}, err
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: uuid,
	}, nil
}
