package payment

import "errors"

var (
	// ErrPaymentNotFound indicates that payment was not found in storage
	ErrPaymentNotFound = errors.New("payment not found")
	// ErrPaymentRequired indicates that payment object is required
	ErrPaymentRequired = errors.New("payment is required")
	// ErrTxIDRequired indicates that transaction ID is required
	ErrTxIDRequired = errors.New("transaction id required")
	// ErrUserUUIDRequired indicates that user UUID is required
	ErrUserUUIDRequired = errors.New("user uuid is required")
	// ErrOrderUUIDRequired indicates that order UUID is required
	ErrOrderUUIDRequired = errors.New("order uuid is required")
	// ErrPaymentMethodInvalid indicates that payment method is required
	ErrPaymentMethodInvalid = errors.New("payment method is required")
)
