package model

import "errors"

var (
	// ErrOrderNotFound возвращается, когда заказ с заданным UUID не найден в БД.
	ErrOrderNotFound = errors.New("order not found")

	// ErrOrderAlreadyCancelled возвращается, если заказ уже был отменён.
	ErrOrderAlreadyCancelled = errors.New("order already cancelled")

	// ErrOrderAlreadyPaid возвращается, если заказ уже был оплачен.
	ErrOrderAlreadyPaid = errors.New("order already paid")

	// ErrCreateOrderFailed возвращается при неудачной попытке создать заказ.
	ErrCreateOrderFailed = errors.New("failed to create order")

	// ErrUpdateOrderFailed возвращается при ошибке обновления заказа.
	ErrUpdateOrderFailed = errors.New("failed to update order")

	// ErrInvalidOrderUUID возвращается, если переданный UUID некорректен.
	ErrInvalidOrderUUID = errors.New("invalid order uuid")

	// ErrDatabase возвращается при непредвиденной ошибке БД.
	ErrDatabase = errors.New("database error")
)
