package converter

import (
	"github.com/agumiroff/BigTechProject/payment/v1/internal/model"
	repomodel "github.com/agumiroff/BigTechProject/payment/v1/internal/repository/model"
)

func ModelToRepo(m *model.Payment) *repomodel.Payment {
	return &repomodel.Payment{
		UserUuid:      m.UserUuid,
		OrderUuid:     m.OrderUuid,
		PaymentMethod: paymentToRepo(m.PaymentMethod),
	}
}

func RepoToModel(r repomodel.Payment) *model.Payment {
	return &model.Payment{
		UserUuid:      r.UserUuid,
		OrderUuid:     r.OrderUuid,
		PaymentMethod: paymentToModel(r.PaymentMethod),
	}
}

func paymentToRepo(p model.PaymentMethod) repomodel.PaymentMethod {
	switch p {
	case model.CARD:
		return repomodel.CARD
	case model.SBP:
		return repomodel.SBP
	case model.CreditCard:
		return repomodel.CreditCard
	case model.InvestorMoney:
		return repomodel.InvestorMoney
	default:
		return repomodel.CategoryUnspecified
	}
}

func paymentToModel(p repomodel.PaymentMethod) model.PaymentMethod {
	switch p {
	case repomodel.CARD:
		return model.CARD
	case repomodel.SBP:
		return model.SBP
	case repomodel.CreditCard:
		return model.CreditCard
	case repomodel.InvestorMoney:
		return model.InvestorMoney
	default:
		return model.CategoryUnspecified
	}
}
