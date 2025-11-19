package payment

import (
	"context"

	"go.uber.org/zap"

	"github.com/agumiroff/BigTechProject/payment/v1/internal/model"
	"github.com/agumiroff/BigTechProject/payment/v1/internal/repository/payment"
	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
)

func (s *service) PayOrder(ctx context.Context, p *model.Payment) (string, error) {
	if err := validatePayment(ctx, p); err != nil {
		return "", err
	}

	res, err := s.Repo.PayOrder(ctx, p)
	if err != nil {
		return "", err
	}

	return res, nil
}

func validatePayment(ctx context.Context, p *model.Payment) error {
	if p == nil {
		logger.Warn(ctx, "Payment validation failed: payment is nil")
		return payment.ErrPaymentRequired
	}

	if p.OrderUUID == "" {
		logger.Warn(ctx, "Payment validation failed: OrderUUID is empty")
		return payment.ErrOrderUUIDRequired
	}

	switch p.PaymentMethod {
	case model.PaymentMethodCard, model.PaymentMethodSBP, model.PaymentMethodCreditCard, model.PaymentMethodInvestMoney:
		logger.Debug(ctx, "Payment validation successful", zap.String("payment_method", string(p.PaymentMethod)))
		return nil
	default:
		logger.Warn(ctx, "Payment validation failed: Invalid payment method", zap.String("payment_method", string(p.PaymentMethod)))
		return payment.ErrPaymentMethodInvalid
	}
}
