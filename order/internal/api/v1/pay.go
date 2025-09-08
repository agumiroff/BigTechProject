package api

import (
	"context"
	"errors"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	orderV1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (*orderV1.PayOrderResponse, error) {
	if req == nil || req.GetPaymentMethod() == "" {
		return nil, errors.New("validation error: missing required fields")
	}

	r, err := a.service.PayOrder(ctx, &model.PayOrderRequest{
		OrderUUID:     params.OrderUUID.String(),
		PaymentMethod: model.PaymentMethod(req.GetPaymentMethod()),
	})
	if err != nil {
		return nil, err
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: r.TransactionUUID,
	}, nil
}
