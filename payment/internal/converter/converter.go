package converter

import (
	"github.com/agumiroff/BigTechProject/payment/v1/internal/model"
	PayV1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

func PaymentToProto(p *model.Payment) *PayV1.Payment {
	if p == nil {
		return nil
	}
	return &PayV1.Payment{
		OrderUuid:     p.OrderUUID,
		UserUuid:      p.UUID,
		PaymentMethod: paymentMethodToProto(p.PaymentMethod),
	}
}

func PaymentToModel(p *PayV1.Payment) *model.Payment {
	if p == nil {
		return nil
	}
	return &model.Payment{
		OrderUUID:     p.OrderUuid,
		UUID:          p.UserUuid,
		PaymentMethod: paymentMethodToModel(p.PaymentMethod),
	}
}

func paymentMethodToProto(p model.PaymentMethod) PayV1.PaymentMethod {
	switch p {
	case model.PaymentMethodCard:
		return PayV1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return PayV1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCreditCard:
		return PayV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodInvestMoney:
		return PayV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return PayV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

func paymentMethodToModel(p PayV1.PaymentMethod) model.PaymentMethod {
	switch p {
	case PayV1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case PayV1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case PayV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case PayV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestMoney
	default:
		return ""
	}
}
