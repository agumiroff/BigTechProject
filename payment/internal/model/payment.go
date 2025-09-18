package model

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "PENDING"
	PaymentStatusProcessing PaymentStatus = "PROCESSING"
	PaymentStatusCompleted  PaymentStatus = "COMPLETED"
	PaymentStatusFailed     PaymentStatus = "FAILED"
	PaymentStatusRefunded   PaymentStatus = "REFUNDED"
)

type PaymentMethod string

const (
	PaymentMethodCard        PaymentMethod = "CARD"
	PaymentMethodSBP         PaymentMethod = "SBP"
	PaymentMethodCreditCard  PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestMoney PaymentMethod = "INVEST_MONEY"
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
