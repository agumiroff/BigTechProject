package order

import (
	"context"

	"go.uber.org/zap"

	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func (s *service) DeleteOrder(ctx context.Context, uuid string) error {
	if uuid == "" {
		return apperrors.ErrInvalidRequest
	}

	order, _, err := s.Repo.GetOrder(ctx, uuid)
	if err != nil {
		logger.Error(ctx, "Order with uuid does not exist", zap.String("uuid", uuid), zap.Error(err))
		return err
	}

	err = s.Repo.DeleteOrder(ctx, order.OrderUUID)
	if err != nil {
		return err
	}
	logger.Info(ctx, "Order deleted", zap.String("order_uuid", order.OrderUUID))

	return nil
}
