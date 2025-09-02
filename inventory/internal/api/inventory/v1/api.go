package inventory

import (
	"github.com/agumiroff/BigTechProject/inventory/v1/internal/service"
	InvV1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
)

type api struct {
	InvV1.UnimplementedInventoryServiceServer

	service service.InvService
}

func NewAPI(service service.InvService) *api {
	return &api{
		service: service,
	}
}
