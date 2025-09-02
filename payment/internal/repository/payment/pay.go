package payment

import (
	"context"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/agumiroff/BigTechProject/payment/v1/internal/model"
	"github.com/agumiroff/BigTechProject/payment/v1/internal/repository/converter"
)

func (r *repository) PayOrder(ctx context.Context, p *model.Payment) (string, error) {
	if p == nil {
		return "", ErrPaymentRequired
	}

	if p.UserUuid == "" {
		return "", ErrUserUUIDRequired
	}

	if p.OrderUuid == "" {
		return "", ErrOrderUUIDRequired
	}

	if p.PaymentMethod == model.CategoryUnspecified {
		return "", ErrPaymentMethodInvalid
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	uuid := gofakeit.UUID()
	r.storage[uuid] = converter.ModelToRepo(p)

	return uuid, nil
}
