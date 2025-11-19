package order

import (
	"context"

	"go.uber.org/zap"

	"github.com/agumiroff/BigTechProject/order/v1/internal/converter"
	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func (s *service) GetOrder(ctx context.Context, uuid string) (*model.Order, error) {
	if uuid == "" {
		return nil, apperrors.ErrInvalidRequest
	}
	order, parts, err := s.Repo.GetOrder(ctx, uuid)
	if err != nil {
		logger.Error(ctx, "failed to get order", zap.String("uuid", uuid), zap.Error(err))
		return nil, err
	}

	return converter.ToModelOrder(order, parts), nil
}
