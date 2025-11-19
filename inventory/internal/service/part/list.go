package part

import (
	"context"

	"go.uber.org/zap"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
	rConverter "github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/converter"
	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
)

func (s *service) ListParts(ctx context.Context, f *model.PartsFilter) (res []*model.Part, err error) {
	part, err := s.Repo.ListParts(ctx, rConverter.FilterToRepo(f))
	if err != nil {
		logger.Error(ctx, "part not found", zap.Error(err))
		return []*model.Part{}, err
	}

	return part, nil
}
