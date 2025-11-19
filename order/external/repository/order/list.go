package order

import (
	"context"

	"go.uber.org/zap"

	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
	inventoryv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
)

func (r *repository) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	res, err := r.invClient.ListParts(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to get list of parts", zap.Error(err))
		return &inventoryv1.ListPartsResponse{}, err
	}

	return res, nil
}
