package model

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
	OrderUUID     string        `json:"order_uuid"`
	UUID          string        `json:"uuid"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	Amount        float64       `json:"amount"`
}

type PaymentFilter struct {
	UUIDs      []string        `json:"uuids,omitempty"`
	OrderUUIDs []string        `json:"order_uuids,omitempty"`
	Statuses   []PaymentStatus `json:"statuses,omitempty"`
}
