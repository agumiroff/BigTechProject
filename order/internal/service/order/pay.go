package order

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/agumiroff/BigTechProject/order/v1/internal/converter"
	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	repomodel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	ordererrors "github.com/agumiroff/BigTechProject/order/v1/internal/service/errors"
	paymentv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
	"github.com/agumiroff/BigTechProject/shared/v1/apperrors"
)

func (s *service) PayOrder(ctx context.Context, req *model.PayOrderRequest) (*model.PayOrderResponse, error) {
	var err error
	var order *repomodel.Order
	var txid *paymentv1.PayOrderResponse

	if req.OrderUUID == "" {
		return nil, apperrors.ErrInvalidRequest
	}

	// Get the order to verify status and get UserUUID
	if order, err = s.Repo.Get(ctx, req.OrderUUID); err != nil {
		log.Printf("failed to get order #%v\n %v", req.OrderUUID, err)
		return nil, err
	}

	if order.Status == repomodel.OrderStatusPAID {
		return nil, ordererrors.ErrOrderPaid
	}

	if order.Status == repomodel.OrderStatusCANCELLED {
		return nil, ordererrors.ErrOrderCancelled
	}

	// Process payment
	if txid, err = s.ExRepo.PayOrder(ctx, &paymentv1.PayOrderRequest{
		Payment: &paymentv1.Payment{
			OrderUuid:     req.OrderUUID,
			UserUuid:      order.UserUUID,
			PaymentMethod: converter.ToProtoPaymentMethod(&req.PaymentMethod),
		},
	}); err != nil {
		log.Printf("Failed to pay order %v", err)
		return nil, err
	}

	// Update order status
	order.PaymentMethod = repomodel.PaymentMethod(req.PaymentMethod)
	order.TransactionUUID = txid.TransactionUuid
	order.Status = repomodel.OrderStatusPAID

	if err = s.Repo.UpdateOrder(ctx, converter.ToModelOrder(order)); err != nil {
		log.Printf("failed to update order #%v\n %v", req.OrderUUID, err)
		return nil, err
	}

	uuid, err := uuid.Parse(txid.TransactionUuid)
	if err != nil {
		return nil, err
	}

	resp := &model.PayOrderResponse{
		TransactionUUID: uuid,
	}

	return resp, nil
}
