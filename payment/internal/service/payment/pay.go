package payment

import (
	"context"
	"log"

	"github.com/agumiroff/BigTechProject/payment/v1/internal/model"
	"github.com/agumiroff/BigTechProject/payment/v1/internal/repository/payment"
)

func (s *service) PayOrder(ctx context.Context, p *model.Payment) (string, error) {
	if err := validatePayment(p); err != nil {
		return "", err
	}

	res, err := s.Repo.PayOrder(ctx, p)
	if err != nil {
		return "", err
	}

	return res, nil
}

func validatePayment(p *model.Payment) error {
	if p == nil {
		log.Printf("Payment validation failed: payment is nil")
		return payment.ErrPaymentRequired
	}

	if p.OrderUUID == "" {
		log.Printf("Payment validation failed: OrderUUID is empty")
		return payment.ErrOrderUUIDRequired
	}

	switch p.PaymentMethod {
	case model.PaymentMethodCard, model.PaymentMethodSBP, model.PaymentMethodCreditCard, model.PaymentMethodInvestMoney:
		log.Printf("Payment validation successful: Valid payment method: %s", p.PaymentMethod)
		return nil
	default:
		log.Printf("Payment validation failed: Invalid payment method: %s", p.PaymentMethod)
		return payment.ErrPaymentMethodInvalid
	}
}
