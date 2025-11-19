package payment

import (
	"context"

	"go.uber.org/zap"

	"github.com/agumiroff/BigTechProject/payment/v1/internal/model"
	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
)

func (s *service) GetPayment(ctx context.Context, uuid string) (*model.Payment, error) {
	payment, err := s.Repo.GetPayment(ctx, uuid)
	if err != nil {
		logger.Error(ctx, "failed to get payment", zap.Error(err), zap.String("uuid", uuid))
		return nil, err
	}

	return payment, nil
}
