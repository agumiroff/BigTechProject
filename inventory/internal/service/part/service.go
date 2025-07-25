package part

import (
	r "github.com/agumiroff/BigTechProject/inventory/v1/internal/repository"
	def "github.com/agumiroff/BigTechProject/inventory/v1/internal/service"
)

var _ def.InvService = (*service)(nil)

type service struct {
	Repo r.InvRepository
}

func NewService(r r.InvRepository) *service {
	return &service{
		Repo: r,
	}
}
