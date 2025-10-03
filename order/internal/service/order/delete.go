package order

import (
	"context"
	"log"

	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func (s *service) DeleteOrder(ctx context.Context, uuid string) error {
	if uuid == "" {
		return apperrors.ErrInvalidRequest
	}

	order, _, err := s.Repo.GetOrder(ctx, uuid)
	if err != nil {
		log.Printf("Order with uuid does not exist %v: %v", uuid, err)
		return err
	}

	err = s.Repo.DeleteOrder(ctx, order.OrderUUID)
	if err != nil {
		return err
	}
	log.Printf("Order deleted")

	return nil
}
