package order

import (
	"github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
)

func (r *repository) Get(uuid string) (*model.Order, error) {
	if uuid == "" {
		return nil, model.ErrInvalidOrderUUID
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.storage[uuid]
	if !exists {
		return nil, model.ErrOrderNotFound
	}

	return order, nil
}
