package api

import (
	"context"

	"go.uber.org/zap"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/converter"
	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
	invV1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *invV1.ListPartsRequest) (*invV1.ListPartsResponse, error) {
	filter := req.GetFilter()
	if filter == nil {
		filter = &invV1.PartsFilter{}
	}

	list, err := a.service.ListParts(ctx, converter.ToFilterModel(filter))
	if err != nil {
		logger.Error(ctx, "failed to list parts", zap.Error(err))
		return &invV1.ListPartsResponse{
			Parts: []*invV1.Part{},
		}, apperrors.Map(err)
	}

	return &invV1.ListPartsResponse{
		Parts: converter.ToProtoModels(list),
	}, nil
}
