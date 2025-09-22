package model

type Payment struct {
	UserUuid      string
	OrderUuid     string
	PaymentMethod PaymentMethod
}

type PaymentMethod int32

const (
	CategoryUnspecified PaymentMethod = 0
	CARD                PaymentMethod = 1
	SBP                 PaymentMethod = 2
	CreditCard          PaymentMethod = 3
	InvestorMoney       PaymentMethod = 4
)
