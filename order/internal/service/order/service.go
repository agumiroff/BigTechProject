package order

import (
	exRepo "github.com/agumiroff/BigTechProject/order/v1/external/repository"
	repo "github.com/agumiroff/BigTechProject/order/v1/internal/repository"
)

type service struct {
	Repo   repo.OrderRepository
	ExRepo exRepo.OrderRepository
}

func NewService(repo repo.OrderRepository, exRepo exRepo.OrderRepository) *service {
	return &service{
		ExRepo: exRepo,
		Repo:   repo,
	}
}
