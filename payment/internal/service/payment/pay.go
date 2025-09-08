package payment

import (
	"context"
	"log"

	"github.com/agumiroff/BigTechProject/payment/v1/internal/model"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func (s *service) PayOrder(ctx context.Context, p *model.Payment) (string, error) {
	if err := validatePayment(p); err != nil {
		return "", err
	}

	res, err := s.Repo.PayOrder(ctx, p)
	if err != nil {
		log.Printf("failed to pay order %v", err)
		return "", err
	}

	return res, nil
}

func validatePayment(p *model.Payment) error {
	if p == nil {
		return apperrors.ErrInvalidRequest
	}

	if p.OrderUuid == "" || p.UserUuid == "" {
		return apperrors.ErrInvalidRequest
	}

	switch p.PaymentMethod {
	case model.CARD, model.SBP, model.CreditCard, model.InvestorMoney:
		return nil
	default:
		return apperrors.ErrInvalidRequest
	}
}
