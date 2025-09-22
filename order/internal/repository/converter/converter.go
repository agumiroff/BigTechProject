package converter

import (
	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	repomodel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
)

func ModelPayToRepo(m *model.PayOrderRequest) *repomodel.PayOrderRequest {
	return &repomodel.PayOrderRequest{
		OrderUUID:     m.OrderUUID,
		PaymentMethod: repomodel.PaymentMethod(m.PaymentMethod),
	}
}

func ModelOrderToRepo(m *model.Order) *repomodel.Order {
	return &repomodel.Order{
		OrderUUID:       m.OrderUUID,
		Status:          repomodel.OrderStatus(m.Status),
		UserUUID:        m.UserUUID,
		PartUUIDs:       m.PartUUIDs,
		TotalPrice:      m.TotalPrice,
		TransactionUUID: m.TransactionUUID,
		PaymentMethod:   repomodel.PaymentMethod(m.PaymentMethod),
	}
}
