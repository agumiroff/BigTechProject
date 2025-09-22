package order

import (
	"context"
	"log"

	paymentv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

func (r *repository) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	res, err := r.payClient.PayOrder(ctx, req)
	if err != nil {
		log.Printf("Failed to pay order: %v", err)
		return &paymentv1.PayOrderResponse{}, err
	}

	return res, nil
}
