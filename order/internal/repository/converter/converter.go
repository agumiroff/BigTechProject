package converter

import (
	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	rModel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
)

func ModelPayToRepo(m *model.PayOrderRequest) *rModel.PayOrderRequest {
	return &rModel.PayOrderRequest{
		OrderUUID:     m.OrderUUID,
		PaymentMethod: rModel.PaymentMethod(m.PaymentMethod),
	}
}

func ModelOrderToRepo(m *model.Order) *rModel.Order {
	return &rModel.Order{
		OrderUUID:       m.OrderUUID,
		Status:          rModel.OrderStatus(m.Status),
		UserUUID:        m.UserUUID,
		PartUUIDs:       m.PartUUIDs,
		TotalPrice:      m.TotalPrice,
		TransactionUUID: m.TransactionUUID,
		PaymentMethod:   rModel.PaymentMethod(m.PaymentMethod),
	}
}
