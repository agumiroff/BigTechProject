package order

import (
	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	rModel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
)

func (r *repository) CreateOrder(req *model.Order) *model.CreateOrderResponse {
	if req == nil {
		return nil
	}

	if req.OrderUUID == "" {
		return nil
	}

	if req.UserUUID == "" || len(req.PartUUIDs) == 0 {
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[req.OrderUUID] = &rModel.Order{
		UserUUID:   req.UserUUID,
		OrderUUID:  req.OrderUUID,
		PartUUIDs:  req.PartUUIDs,
		TotalPrice: req.TotalPrice,
		Status:     rModel.OrderStatus(req.Status),
	}

	return &model.CreateOrderResponse{
		OrderUUID:  req.OrderUUID,
		TotalPrice: req.TotalPrice,
	}
}
