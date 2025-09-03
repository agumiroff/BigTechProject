package order

import (
	"context"
	"log"

	"github.com/agumiroff/BigTechProject/order/v1/internal/converter"
	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	rModel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	paymentv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

func (s *service) PayOrder(ctx context.Context, req *model.PayOrderRequest) (*model.PayOrderResponse, error) {
	txid, err := s.ExRepo.PayOrder(ctx, &paymentv1.PayOrderRequest{
		Payment: &paymentv1.Payment{
			OrderUuid:     req.OrderUUID,
			PaymentMethod: converter.ToProtoPaymentMethod(&req.PaymentMethod),
		},
	})
	if err != nil {
		log.Printf("Failed to pay order %v", err)
		return &model.PayOrderResponse{}, err
	}

	order, err := s.Repo.Get(ctx, req.OrderUUID)
	if err != nil {
		log.Printf("failed to get order #%v\n %v", req.OrderUUID, err)
		return nil, err
	}

	order.PaymentMethod = rModel.PaymentMethod(req.PaymentMethod)
	order.TransactionUUID = txid.TransactionUuid
	order.Status = rModel.OrderStatusPAID

	err = s.Repo.UpdateOrder(ctx, converter.ToModelOrder(order))
	if err != nil {
		log.Printf("failed to update order #%v\n %v", req.OrderUUID, err)
		return nil, err
	}

	return &model.PayOrderResponse{
		TransactionUUID: txid.TransactionUuid,
	}, nil
}
