package order

import (
	"context"

	"github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func (r *repository) Get(ctx context.Context, uuid string) (*model.Order, error) {
	if uuid == "" {
		return nil, apperrors.ErrInvalidRequest
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.storage[uuid]
	if !exists {
		return nil, apperrors.ErrNotFound
	}

	return order, nil
}
