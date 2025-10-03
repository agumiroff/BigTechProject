package model

type CreateOrderRequest struct {
	UserUUID   string
	PartUuids  []string
	TotalPrice float64
}

type CreateOrderResponse struct {
	OrderUUID  string
	TotalPrice float64
}
