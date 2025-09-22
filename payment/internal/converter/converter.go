package converter

import (
	"github.com/agumiroff/BigTechProject/payment/v1/internal/model"
	PayV1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

func PaymentToProto(p *model.Payment) *PayV1.Payment {
	return &PayV1.Payment{
		UserUuid:      p.UserUuid,
		OrderUuid:     p.OrderUuid,
		PaymentMethod: paymentMethodToProto(p.PaymentMethod),
	}
}

func PaymentToModel(p *PayV1.Payment) *model.Payment {
	if p == nil {
		return nil
	}
	return &model.Payment{
		UserUuid:      p.UserUuid,
		OrderUuid:     p.OrderUuid,
		PaymentMethod: paymentMethodToModel(p.PaymentMethod),
	}
}

func paymentMethodToProto(p model.PaymentMethod) PayV1.PaymentMethod {
	switch p {
	case model.CARD:
		return PayV1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.SBP:
		return PayV1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.CreditCard:
		return PayV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.InvestorMoney:
		return PayV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return PayV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

func paymentMethodToModel(p PayV1.PaymentMethod) model.PaymentMethod {
	switch p {
	case PayV1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.CARD
	case PayV1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.SBP
	case PayV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.CreditCard
	case PayV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.InvestorMoney
	default:
		return model.CategoryUnspecified
	}
}
