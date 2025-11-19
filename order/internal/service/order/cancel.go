package order

import (
	"context"

	"go.uber.org/zap"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	repoConverter "github.com/agumiroff/BigTechProject/order/v1/internal/repository/converter"
	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func (s *service) CancelOrder(ctx context.Context, uuid string) error {
	if uuid == "" {
		return apperrors.ErrInvalidRequest
	}

	order, parts, err := s.Repo.GetOrder(ctx, uuid)
	if err != nil {
		logger.Error(ctx, "Order with uuid does not exist", zap.String("uuid", uuid), zap.Error(err))
		return err
	}

	modelOrder := repoConverter.RepoOrderToModel(order, parts)

	if modelOrder.Status == model.OrderStatusCANCELLED {
		return apperrors.ErrForbidden
	}

	if modelOrder.Status == model.OrderStatusPAID {
		return apperrors.ErrForbidden
	}

	err = s.Repo.CancelOrder(ctx, order.OrderUUID)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Order cancelled", zap.String("order_uuid", order.OrderUUID))

	return nil
}
