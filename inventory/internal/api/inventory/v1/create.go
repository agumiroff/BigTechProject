package api

import (
	"context"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/converter"
	inventoryv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
)

func (a *api) CreatePart(ctx context.Context, req *inventoryv1.CreatePartRequest) (*inventoryv1.CreatePartResponse, error) {
	res, err := a.service.CreatePart(ctx, converter.ToModelPart(req.GetPart()))
	if err != nil {
		return &inventoryv1.CreatePartResponse{}, err
	}

	return &inventoryv1.CreatePartResponse{
		Uuid: res,
	}, nil
}
