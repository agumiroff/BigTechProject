package order

import (
	"sync"

	"github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
)

type repository struct {
	mu sync.RWMutex

	storage map[string]*model.Order
}

func NewRepository() *repository {
	return &repository{
		storage: make(map[string]*model.Order),
	}
}
