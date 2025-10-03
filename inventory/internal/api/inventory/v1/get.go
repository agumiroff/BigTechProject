package inventory

import (
	"context"
	"log"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/converter"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
	inventoryv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (res *inventoryv1.GetPartResponse, err error) {
	part, err := a.service.GetPart(ctx, req.GetUuid())
	if err != nil {
		log.Printf("failed to get part: %v", err)
		return &inventoryv1.GetPartResponse{}, apperrors.Map(err)
	}
	return &inventoryv1.GetPartResponse{
		Part: converter.ToProtoPart(part),
	}, nil
}
