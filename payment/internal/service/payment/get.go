package payment

import (
	"context"
	"log"

	"github.com/agumiroff/BigTechProject/payment/v1/internal/model"
)

func (s *service) GetPayment(ctx context.Context, uuid string) (*model.Payment, error) {
	payment, err := s.Repo.GetPayment(ctx, uuid)
	if err != nil {
		log.Printf("failed to get payment: %v", err)
		return nil, err
	}

	return payment, nil
}
