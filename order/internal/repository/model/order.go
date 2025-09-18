package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

type OrderStatus string

const (
	OrderStatusPENDINGPAYMENT OrderStatus = "PENDING_PAYMENT"
	OrderStatusPAID           OrderStatus = "PAID"
	OrderStatusCANCELLED      OrderStatus = "CANCELLED"
)

type PaymentMethod string

const (
	PaymentMethodUNKNOWN       PaymentMethod = "UNKNOWN"
	PaymentMethodCARD          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCREDITCARD    PaymentMethod = "CREDIT_CARD"
	PaymentMethodINVESTORMONEY PaymentMethod = "INVESTOR_MONEY"
)

type Order struct {
	OrderUUID       string
	UserUUID        string
	PartUUIDs       []string
	TotalPrice      float64
	TransactionUUID string
	PaymentMethod   PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type OrderRow struct {
	OrderUUID       string          `db:"order_uuid"`
	UserUUID        string          `db:"user_uuid"`
	PartUUIDs       json.RawMessage `db:"part_uuids"`
	TotalPrice      float64         `db:"total_price"`
	TransactionUUID sql.NullString  `db:"transaction_uuid"`
	PaymentMethod   sql.NullString  `db:"payment_method"`
	Status          string          `db:"status"`
	CreatedAt       time.Time       `db:"created_at"`
	UpdatedAt       sql.NullTime    `db:"updated_at"`
}
