package order

import (
	"context"

	repomodel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func (r *repository) DeleteOrder(ctx context.Context, uuid string) error {
	if uuid == "" {
		return apperrors.ErrInvalidRequest
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.storage[uuid]
	if !exists {
		return apperrors.ErrNotFound
	}

	if existing.Status == repomodel.OrderStatusCANCELLED {
		return apperrors.ErrForbidden
	}

	delete(r.storage, uuid)
	return nil
}
