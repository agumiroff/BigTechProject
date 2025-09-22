package model

type CreateOrderRequest struct {
	UserUuid   string
	PartUuids  []string
	TotalPrice float64
}

type CreateOrderResponse struct {
	OrderUuid  string
	TotalPrice float64
}
