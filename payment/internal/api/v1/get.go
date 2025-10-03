package v1

import (
	"context"

	"github.com/agumiroff/BigTechProject/payment/v1/internal/converter"
	paymentv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

func (a *API) GetPayment(ctx context.Context, uuid string) (*paymentv1.Payment, error) {
	payment, err := a.service.GetPayment(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return converter.PaymentToProto(payment), nil
}
