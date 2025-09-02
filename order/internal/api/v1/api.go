package api

import (
	"github.com/agumiroff/BigTechProject/order/v1/internal/service"
)

// api is a test implementation of the order API
type api struct {
	service service.OrderService
}

// NewAPI creates a new test API instance
func NewAPI(service service.OrderService) *api {
	return &api{
		service: service,
	}
}
