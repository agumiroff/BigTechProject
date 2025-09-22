package order

import (
	"google.golang.org/grpc"

	inventoryv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

type repository struct {
	invClient inventoryv1.InventoryServiceClient
	payClient paymentv1.PaymentServiceClient
}

func NewRepository(invConn, payConn grpc.ClientConnInterface) *repository {
	return &repository{
		invClient: inventoryv1.NewInventoryServiceClient(invConn),
		payClient: paymentv1.NewPaymentServiceClient(payConn),
	}
}
