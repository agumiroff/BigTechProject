package ordererrors

import "errors"

// Common errors shared across all services
var (
	ErrOrderPaid      = errors.New("order already paid")
	ErrOrderCancelled = errors.New("order cancelled")
)
