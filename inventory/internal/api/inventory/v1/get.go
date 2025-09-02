package inventory

import (
	"context"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/converter"
	inventoryv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (res *inventoryv1.GetPartResponse, err error) {
	part, err := a.service.GetPart(ctx, req.GetUuid())
	if err != nil {
		return &inventoryv1.GetPartResponse{}, err
	}
	return &inventoryv1.GetPartResponse{
		Part: converter.ModelToProto(part),
	}, nil
}
