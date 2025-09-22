package order

import (
	"context"

	"github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func (r *repository) CancelOrder(ctx context.Context, uuid string) error {
	if uuid == "" {
		return apperrors.ErrInvalidRequest
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	order, exists := r.storage[uuid]
	if !exists {
		return apperrors.ErrNotFound
	}

	if order.Status == model.OrderStatusCANCELLED {
		return apperrors.ErrForbidden
	}

	if order.Status == model.OrderStatusPAID {
		return apperrors.ErrForbidden
	}

	order.Status = model.OrderStatusCANCELLED
	r.storage[uuid] = order

	return nil
}
