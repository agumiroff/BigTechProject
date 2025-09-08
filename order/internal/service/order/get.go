package order

import (
	"context"
	"log"

	"github.com/agumiroff/BigTechProject/order/v1/internal/converter"
	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	"github.com/agumiroff/BigTechProject/shared/v1/apperrors"
)

func (s *service) GetOrder(ctx context.Context, uuid string) (*model.Order, error) {
	if uuid == "" {
		return nil, apperrors.ErrInvalidRequest
	}
	order, err := s.Repo.Get(ctx, uuid)
	if err != nil {
		log.Printf("failed to get order %v: %v", uuid, err)
		return nil, err
	}

	return converter.ToModelOrder(order), nil
}
