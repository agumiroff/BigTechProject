package order

import (
	"context"
	"log"
)

func (s *service) CancelOrder(ctx context.Context, uuid string) error {
	order, err := s.Repo.Get(ctx, uuid)
	if err != nil {
		log.Printf("Order with uuid does not exist %v: %v", uuid, err)
		return err
	}

	err = s.Repo.DeleteOrder(ctx, order.OrderUUID)
	if err != nil {
		return err
	}
	log.Printf("Order cancelled")

	return nil
}
