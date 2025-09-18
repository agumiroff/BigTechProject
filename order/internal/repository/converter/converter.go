package converter

import (
	"database/sql"
	"encoding/json"
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

func ModelOrderToRepo(m *model.Order) *repomodel.OrderRow {
	// Serialize partUUIDs to JSON
	partUUIDsJSON, _ := json.Marshal(m.PartUUIDs)

	now := time.Now()
	return &repomodel.OrderRow{
		OrderUUID:  m.OrderUUID,
		UserUUID:   m.UserUUID,
		PartUUIDs:  partUUIDsJSON,
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
		UpdatedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
	}
}

func RepoOrderToModel(r *repomodel.OrderRow) *model.Order {
	var partUUIDs []string
	_ = json.Unmarshal(r.PartUUIDs, &partUUIDs)

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
