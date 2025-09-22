package order

import (
	"context"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	repomodel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func (r *repository) CreateOrder(ctx context.Context, req *model.Order) (*model.CreateOrderResponse, error) {
	if req == nil {
		return nil, apperrors.ErrInvalidRequest
	}

	if req.OrderUUID == "" {
		return nil, apperrors.ErrInvalidRequest
	}

	if req.UserUUID == "" || len(req.PartUUIDs) == 0 {
		return nil, apperrors.ErrInvalidRequest
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.storage[req.OrderUUID]; exists {
		return nil, apperrors.ErrAlreadyExists
	}

	r.storage[req.OrderUUID] = &repomodel.Order{
		UserUUID:   req.UserUUID,
		OrderUUID:  req.OrderUUID,
		PartUUIDs:  req.PartUUIDs,
		TotalPrice: req.TotalPrice,
		Status:     repomodel.OrderStatus(req.Status),
	}

	return &model.CreateOrderResponse{
		OrderUUID:  req.OrderUUID,
		TotalPrice: req.TotalPrice,
	}, nil
}
