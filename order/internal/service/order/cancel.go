package order

import (
	"context"
	"log"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	repoConverter "github.com/agumiroff/BigTechProject/order/v1/internal/repository/converter"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func (s *service) CancelOrder(ctx context.Context, uuid string) error {
	if uuid == "" {
		return apperrors.ErrInvalidRequest
	}

	order, parts, err := s.Repo.GetOrder(ctx, uuid)
	if err != nil {
		log.Printf("Order with uuid does not exist %v: %v", uuid, err)
		return err
	}

	modelOrder := repoConverter.RepoOrderToModel(order, parts)

	if modelOrder.Status == model.OrderStatusCANCELLED {
		return apperrors.ErrForbidden
	}

	if modelOrder.Status == model.OrderStatusPAID {
		return apperrors.ErrForbidden
	}

	err = s.Repo.CancelOrder(ctx, order.OrderUUID)
	if err != nil {
		return err
	}

	log.Printf("Order cancelled")

	return nil
}
