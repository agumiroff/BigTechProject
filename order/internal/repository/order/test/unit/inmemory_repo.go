package test

import (
	"context"
	"sync"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	repomodel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

type inMemoryRepo struct {
	mu      sync.RWMutex
	storage map[string]*repomodel.Order
}

func NewInmemoryRepo() *inMemoryRepo {
	return &inMemoryRepo{
		mu:      sync.RWMutex{},
		storage: make(map[string]*repomodel.Order),
	}
}

func (r *inMemoryRepo) CreateOrder(ctx context.Context, order *model.Order) (*model.CreateOrderResponse, error) {
	if order == nil {
		return nil, apperrors.ErrInvalidRequest
	}

	if order.OrderUUID == "" {
		return nil, apperrors.ErrInvalidRequest
	}

	if order.UserUUID == "" || len(order.PartUUIDs) == 0 {
		return nil, apperrors.ErrInvalidRequest
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.storage[order.OrderUUID]; exists {
		return nil, apperrors.ErrAlreadyExists
	}

	r.storage[order.OrderUUID] = &repomodel.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		Status:          repomodel.OrderStatus(order.Status),
		PaymentMethod:   repomodel.PaymentMethod(order.PaymentMethod),
		TransactionUUID: order.TransactionUUID,
	}

	return &model.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}

func (r *inMemoryRepo) Get(ctx context.Context, uuid string) (*repomodel.Order, error) {
	if uuid == "" {
		return nil, apperrors.ErrInvalidRequest
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.storage[uuid]
	if !exists {
		return nil, apperrors.ErrNotFound
	}

	return order, nil
}

func (r *inMemoryRepo) UpdateOrder(ctx context.Context, order *model.Order) error {
	if order == nil {
		return apperrors.ErrInvalidRequest
	}

	if order.OrderUUID == "" {
		return apperrors.ErrInvalidRequest
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.storage[order.OrderUUID]
	if !exists {
		return apperrors.ErrNotFound
	}

	if existing.Status == repomodel.OrderStatusCANCELLED {
		return apperrors.ErrForbidden
	}

	if existing.Status == repomodel.OrderStatusPAID &&
		order.Status != model.OrderStatusCANCELLED {
		return apperrors.ErrForbidden
	}

	r.storage[order.OrderUUID] = &repomodel.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		Status:          repomodel.OrderStatus(order.Status),
		PaymentMethod:   repomodel.PaymentMethod(order.PaymentMethod),
		TransactionUUID: order.TransactionUUID,
	}

	return nil
}

func (r *inMemoryRepo) DeleteOrder(ctx context.Context, uuid string) error {
	if uuid == "" {
		return apperrors.ErrInvalidRequest
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.storage[uuid]
	if !exists {
		return apperrors.ErrNotFound
	}

	if existing.Status == repomodel.OrderStatusCANCELLED {
		return apperrors.ErrForbidden
	}

	delete(r.storage, uuid)
	return nil
}

func (r *inMemoryRepo) CancelOrder(ctx context.Context, uuid string) error {
	if uuid == "" {
		return apperrors.ErrInvalidRequest
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	order, exists := r.storage[uuid]
	if !exists {
		return apperrors.ErrNotFound
	}

	if order.Status == repomodel.OrderStatusCANCELLED {
		return apperrors.ErrForbidden
	}

	if order.Status == repomodel.OrderStatusPAID {
		return apperrors.ErrForbidden
	}

	order.Status = repomodel.OrderStatusCANCELLED
	r.storage[uuid] = order

	return nil
}
