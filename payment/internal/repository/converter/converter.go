package converter

import (
	"github.com/agumiroff/BigTechProject/payment/v1/internal/model"
	rModel "github.com/agumiroff/BigTechProject/payment/v1/internal/repository/model"
)

func ModelToRepo(m *model.Payment) *rModel.Payment {
	return &rModel.Payment{
		UserUuid:      m.UserUuid,
		OrderUuid:     m.OrderUuid,
		PaymentMethod: paymentToRepo(m.PaymentMethod),
	}
}

func RepoToModel(r rModel.Payment) *model.Payment {
	return &model.Payment{
		UserUuid:      r.UserUuid,
		OrderUuid:     r.OrderUuid,
		PaymentMethod: paymentToModel(r.PaymentMethod),
	}
}

func paymentToRepo(p model.PaymentMethod) rModel.PaymentMethod {
	switch p {
	case model.CARD:
		return rModel.CARD
	case model.SBP:
		return rModel.SBP
	case model.CreditCard:
		return rModel.CreditCard
	case model.InvestorMoney:
		return rModel.InvestorMoney
	default:
		return rModel.CategoryUnspecified
	}
}

func paymentToModel(p rModel.PaymentMethod) model.PaymentMethod {
	switch p {
	case rModel.CARD:
		return model.CARD
	case rModel.SBP:
		return model.SBP
	case rModel.CreditCard:
		return model.CreditCard
	case rModel.InvestorMoney:
		return model.InvestorMoney
	default:
		return model.CategoryUnspecified
	}
}
