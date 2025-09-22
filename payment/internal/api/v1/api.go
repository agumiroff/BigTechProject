// Package v1 implements gRPC API handlers for payment service
package v1

import (
	"github.com/agumiroff/BigTechProject/payment/v1/internal/service"
	paymentv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

// API handles gRPC requests for payment service
type API struct {
	paymentv1.UnimplementedPaymentServiceServer

	service service.PaymentService
}

// NewAPI creates a new payment API handler
func NewAPI(service service.PaymentService) *API {
	return &API{
		service: service,
	}
}
