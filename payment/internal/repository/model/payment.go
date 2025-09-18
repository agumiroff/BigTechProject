package model

import "time"

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusCompleted  PaymentStatus = "completed"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusRefunded   PaymentStatus = "refunded"
)

type PaymentMethod string

const (
	PaymentMethodCard        PaymentMethod = "card"
	PaymentMethodSBP         PaymentMethod = "sbp"
	PaymentMethodCreditCard  PaymentMethod = "credit_card"
	PaymentMethodInvestMoney PaymentMethod = "invest_money"
)

type Payment struct {
	UUID          string        `json:"uuid" bson:"uuid"`
	OrderUUID     string        `json:"order_uuid" bson:"order_uuid"`
	Status        PaymentStatus `json:"status" bson:"status"`
	Amount        float64       `json:"amount" bson:"amount"`
	PaymentMethod PaymentMethod `json:"method" bson:"method"`
	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" bson:"updated_at"`
}
