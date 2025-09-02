package payment

import (
	"context"
	"log"

	"github.com/agumiroff/BigTechProject/payment/v1/internal/model"
)

func (s *service) PayOrder(ctx context.Context, p *model.Payment) (string, error) {
	res, err := s.Repo.PayOrder(ctx, p)
	if err != nil {
		log.Printf("failed to pay order %v", err)
		return "", err
	}

	return res, nil
}
