package payment

import (
	"context"

	"github.com/agumiroff/BigTechProject/payment/v1/internal/model"
	"github.com/agumiroff/BigTechProject/payment/v1/internal/repository/converter"
)

func (r *repository) GetPayment(ctx context.Context, uuid string) (*model.Payment, error) {
	if uuid == "" {
		return nil, ErrTxIDRequired
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	payment, exists := r.storage[uuid]
	if !exists {
		return nil, ErrPaymentNotFound
	}

	return converter.RepoToModel(*payment), nil
}
