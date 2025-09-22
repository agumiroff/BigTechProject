package order

import (
	"context"
	"log"

	inventoryv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
)

func (r *repository) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	res, err := r.invClient.ListParts(ctx, req)
	if err != nil {
		log.Printf("Failed to get list of parts %v", err)
		return &inventoryv1.ListPartsResponse{}, err
	}

	return res, nil
}
