package order

import (
	"context"
	"log"
)

func (s *service) DeleteOrder(ctx context.Context, uuid string) error {
	order, err := s.Repo.Get(uuid)
	if err != nil {
		log.Printf("failed to find order: %v", err)
		return err
	}

	err = s.Repo.DeleteOrder(order.OrderUUID)
	if err != nil {
		log.Printf("failed to delete order: %v", err)
		return err
	}

	return nil
}
