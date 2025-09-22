package order

import (
	"sync"

	repomodel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
)

type repository struct {
	mu sync.RWMutex

	storage map[string]*repomodel.Order
}

func NewRepository() *repository {
	return &repository{
		storage: make(map[string]*repomodel.Order),
	}
}
