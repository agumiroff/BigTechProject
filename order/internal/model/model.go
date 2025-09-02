package model

type CreateOrderRequest struct {
	UserUUID  string
	PartUUIDs []string
}

type CreateOrderResponse struct {
	OrderUUID  string
	TotalPrice float64
}
