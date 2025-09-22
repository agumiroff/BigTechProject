package model

type PayOrderRequest struct {
	OrderUUID     string
	UserUUID      string
	PaymentMethod PaymentMethod
}

type PayOrderResponse struct {
	TransactionUUID string
}
