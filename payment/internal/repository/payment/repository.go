package payment

import (
	"sync"

	"github.com/agumiroff/BigTechProject/payment/v1/internal/repository/model"
)

type repository struct {
	mu sync.RWMutex

	storage map[string]*model.Payment
}

func NewRepository() *repository {
	return &repository{
		storage: make(map[string]*model.Payment),
	}
}
