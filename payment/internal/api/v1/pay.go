package v1

import (
	"context"

	"go.uber.org/zap"

	"github.com/agumiroff/BigTechProject/payment/v1/internal/converter"
	"github.com/agumiroff/BigTechProject/payment/v1/internal/repository/payment"
	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
	paymentv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

// PayOrder processes a payment request
func (a *API) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	if req == nil || req.Payment == nil {
		return &paymentv1.PayOrderResponse{}, payment.ErrPaymentRequired
	}
	res, err := a.service.PayOrder(ctx, converter.PaymentToModel(req.Payment))
	if err != nil {
		logger.Error(ctx, "There was a error trying to pay order", zap.Error(err))
		return &paymentv1.PayOrderResponse{}, err
	}

	return &paymentv1.PayOrderResponse{
		TransactionUuid: res,
	}, nil
}
