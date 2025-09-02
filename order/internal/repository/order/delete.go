package order

import "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"

func (r *repository) DeleteOrder(uuid string) error {
	if uuid == "" {
		return model.ErrInvalidOrderUUID
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.storage[uuid]
	if !exists {
		return model.ErrOrderNotFound
	}

	if existing.Status == model.OrderStatusCANCELLED {
		return model.ErrOrderAlreadyCancelled
	}

	delete(r.storage, uuid)
	return nil
}
