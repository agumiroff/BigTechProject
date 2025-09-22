package payment

import (
	r "github.com/agumiroff/BigTechProject/payment/v1/internal/repository"
	def "github.com/agumiroff/BigTechProject/payment/v1/internal/service"
)

var _ def.PaymentService = (*service)(nil)

type service struct {
	Repo r.PaymentRepository
}

func NewService(repo r.PaymentRepository) *service {
	return &service{
		Repo: repo,
	}
}
