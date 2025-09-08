package model

import "github.com/google/uuid"

type PayOrderRequest struct {
	UserUUID      string
	OrderUUID     string
	PaymentMethod PaymentMethod
}

type PayOrderResponse struct {
	TransactionUUID uuid.UUID
}
