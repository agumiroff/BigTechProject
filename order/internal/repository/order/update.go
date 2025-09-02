package order

import (
	"context"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	rModel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
)

func (r *repository) UpdateOrder(ctx context.Context, m *model.Order) error {
	if m == nil {
		return rModel.ErrUpdateOrderFailed
	}

	if m.OrderUUID == "" {
		return rModel.ErrInvalidOrderUUID
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.storage[m.OrderUUID]
	if !exists {
		return rModel.ErrOrderNotFound
	}

	if existing.Status == rModel.OrderStatusCANCELLED {
		return rModel.ErrOrderAlreadyCancelled
	}

	if existing.Status == rModel.OrderStatus(model.OrderStatusPAID) &&
		m.Status != model.OrderStatusCANCELLED {
		return rModel.ErrOrderAlreadyPaid
	}

	r.storage[m.OrderUUID] = &rModel.Order{
		UserUUID:   m.UserUUID,
		OrderUUID:  m.OrderUUID,
		PartUUIDs:  m.PartUUIDs,
		TotalPrice: m.TotalPrice,
		Status:     rModel.OrderStatus(m.Status),
	}

	return nil
}
