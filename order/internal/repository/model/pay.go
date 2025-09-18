package model

type PayOrderRequest struct {
	OrderUUID     string
	PaymentMethod string
}

type PayOrderResponse struct {
	TransactionUUID string
}
