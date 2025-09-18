package order

import (
	"context"
	"log"

	"github.com/agumiroff/BigTechProject/shared/v1/apperrors"
)

func (s *service) CancelOrder(ctx context.Context, uuid string) error {
	if uuid == "" {
		return apperrors.ErrInvalidRequest
	}

	order, err := s.Repo.Get(ctx, uuid)
	if err != nil {
		log.Printf("Order with uuid does not exist %v: %v", uuid, err)
		return err
	}

	err = s.Repo.CancelOrder(ctx, order.OrderUUID)
	if err != nil {
		return err
	}

	log.Printf("Order cancelled")

	return nil
}
