package converter

import (
	"github.com/agumiroff/BigTechProject/payment/v1/internal/model"
	repomodel "github.com/agumiroff/BigTechProject/payment/v1/internal/repository/model"
)

// ModelToRepo converts a domain model payment to a repository model payment
func ModelToRepo(m *model.Payment) *repomodel.Payment {
	if m == nil {
		return nil
	}

	return &repomodel.Payment{
		OrderUUID:     m.OrderUUID,
		PaymentMethod: repomodel.PaymentMethod(m.PaymentMethod),
		Status:        repomodel.PaymentStatusPending,
		Amount:        m.Amount,
	}
}

// RepoToModel converts a repository model payment to a domain model payment
func RepoToModel(r *repomodel.Payment) *model.Payment {
	if r == nil {
		return nil
	}

	return &model.Payment{
		UUID:          r.UUID,
		OrderUUID:     r.OrderUUID,
		PaymentMethod: model.PaymentMethod(r.PaymentMethod),
	}
}
