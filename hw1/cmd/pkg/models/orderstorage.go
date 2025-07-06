package models

import (
	"sync"
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*Order),
	}
}

func (s *OrderStorage) GetOrder(id string) *Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[id]
	if !ok {
		return nil
	}

	return order
}

func (s *OrderStorage) UpdateOrders(id string, order *Order) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[id] = order
}
