package converter

import (
	"database/sql"
	"time"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	repomodel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
)

func ModelPayToRepo(m *model.PayOrderRequest) *repomodel.PayOrderRequest {
	return &repomodel.PayOrderRequest{
		OrderUUID:     m.OrderUUID,
		PaymentMethod: string(m.PaymentMethod),
	}
}

func ModelOrderToRepo(m *model.Order) (*repomodel.OrderRow, []*repomodel.OrderParts) {
	now := time.Now()

	order := &repomodel.OrderRow{
		OrderUUID:  m.OrderUUID,
		UserUUID:   m.UserUUID,
		TotalPrice: m.TotalPrice,
		TransactionUUID: sql.NullString{
			String: m.TransactionUUID,
			Valid:  m.TransactionUUID != "",
		},
		PaymentMethod: sql.NullString{
			String: string(m.PaymentMethod),
			Valid:  string(m.PaymentMethod) != "",
		},
		Status:    string(m.Status),
		CreatedAt: now,
		UpdatedAt: sql.NullTime{Time: now, Valid: true},
	}

	var orderParts []*repomodel.OrderParts
	for _, part := range m.PartUUIDs { // m.Parts — []PartWithQty
		orderParts = append(orderParts, &repomodel.OrderParts{
			OrderUUID: order.OrderUUID,
			PartUUID:  part,
		})
	}

	return order, orderParts
}

func RepoOrderToModel(r *repomodel.OrderRow, partUUIDs []string) *model.Order {

	return &model.Order{
		OrderUUID:       r.OrderUUID,
		UserUUID:        r.UserUUID,
		PartUUIDs:       partUUIDs,
		TotalPrice:      r.TotalPrice,
		TransactionUUID: nullStringToString(r.TransactionUUID),
		PaymentMethod:   model.PaymentMethod(nullStringToString(r.PaymentMethod)),
		Status:          model.OrderStatus(r.Status),
	}
}

func RepoPayToModel(r *repomodel.PayOrderRequest) *model.PayOrderRequest {
	return &model.PayOrderRequest{
		OrderUUID:     r.OrderUUID,
		PaymentMethod: model.PaymentMethod(r.PaymentMethod),
	}
}

// helper
func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
