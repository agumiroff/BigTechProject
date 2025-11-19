package order

import (
	"context"

	"go.uber.org/zap"

	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
	paymentv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

func (r *repository) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	res, err := r.payClient.PayOrder(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to pay order", zap.Error(err))
		return nil, err
	}

	return res, nil
}
