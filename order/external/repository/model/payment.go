package model

type PayOrderRequest struct {
	Payment *Payment
}

type PayOrderResponse struct {
	TransactionUUID string
}

type Payment struct {
	OrderUUID     string
	UserUUID      string
	PaymentMethod PaymentMethod
}

type PaymentMethod string

const (
	PaymentMethodUNKNOWN       PaymentMethod = "UNKNOWN"
	PaymentMethodCARD          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCREDITCARD    PaymentMethod = "CREDIT_CARD"
	PaymentMethodINVESTORMONEY PaymentMethod = "INVESTOR_MONEY"
)
