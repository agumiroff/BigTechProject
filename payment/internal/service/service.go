package service

import (
	"context"

	"github.com/agumiroff/BigTechProject/payment/v1/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, m *model.Payment) (string, error)
	GetPayment(ctx context.Context, uuid string) (*model.Payment, error)
}
