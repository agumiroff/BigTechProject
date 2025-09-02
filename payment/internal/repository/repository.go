package repository

import (
	"context"

	"github.com/agumiroff/BigTechProject/payment/v1/internal/model"
)

type PaymentRepository interface {
	PayOrder(ctx context.Context, p *model.Payment) (string, error)
	GetPayment(ctx context.Context, uuid string) (*model.Payment, error)
}
