package order

import (
	"context"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	rModel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func (r *repository) UpdateOrder(ctx context.Context, m *model.Order) error {
	if m == nil {
		return apperrors.ErrInvalidRequest
	}

	if m.OrderUUID == "" {
		return apperrors.ErrInvalidRequest
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.storage[m.OrderUUID]
	if !exists {
		return apperrors.ErrNotFound
	}

	if existing.Status == rModel.OrderStatusCANCELLED {
		return apperrors.ErrForbidden
	}

	if existing.Status == rModel.OrderStatus(model.OrderStatusPAID) &&
		m.Status != model.OrderStatusCANCELLED {
		return apperrors.ErrForbidden
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
