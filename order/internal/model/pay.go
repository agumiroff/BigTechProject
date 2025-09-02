package model

type PayOrderRequest struct {
	UserUUID      string
	OrderUUID     string
	PaymentMethod PaymentMethod
}

type PayOrderResponse struct {
	TransactionUUID string
}
