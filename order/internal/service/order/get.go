package order

import (
	"context"
	"log"

	"github.com/agumiroff/BigTechProject/order/v1/internal/converter"
	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
)

func (s *service) GetOrder(ctx context.Context, uuid string) (*model.Order, error) {
	order, err := s.Repo.Get(uuid)
	if err != nil {
		log.Printf("failed to get order %v: %v", uuid, err)
		return &model.Order{}, err
	}

	return converter.RepoOrderToModel(order), nil
}
