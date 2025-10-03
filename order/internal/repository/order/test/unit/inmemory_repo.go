package unit

import (
	"context"
	"sync"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	repomodel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

type OrderData struct {
	Order *repomodel.OrderRow
	Parts []string
}

type inMemoryRepo struct {
	mu      sync.RWMutex
	storage map[string]OrderData
}

func NewInmemoryRepo() *inMemoryRepo {
	return &inMemoryRepo{
		mu:      sync.RWMutex{},
		storage: make(map[string]OrderData),
	}
}

func (r *inMemoryRepo) CreateOrder(ctx context.Context, order *repomodel.OrderRow, parts []string) (*model.CreateOrderResponse, error) {
	if order == nil {
		return nil, apperrors.ErrInvalidRequest
	}

	if order.OrderUUID == "" {
		return nil, apperrors.ErrInvalidRequest
	}

	if order.UserUUID == "" || len(parts) == 0 {
		return nil, apperrors.ErrInvalidRequest
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.storage[order.OrderUUID]; exists {
		return nil, apperrors.ErrAlreadyExists
	}

	r.storage[order.OrderUUID] = OrderData{
		Order: order,
		Parts: parts,
	}

	return &model.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}

func (r *inMemoryRepo) GetOrder(ctx context.Context, uuid string) (*repomodel.OrderRow, []string, error) {
	if uuid == "" {
		return nil, nil, apperrors.ErrInvalidRequest
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	orderData, exists := r.storage[uuid]
	if !exists {
		return nil, nil, apperrors.ErrNotFound
	}

	return orderData.Order, orderData.Parts, nil
}

func (r *inMemoryRepo) UpdateOrder(ctx context.Context, order *repomodel.OrderRow) error {
	if order == nil {
		return apperrors.ErrInvalidRequest
	}

	if order.OrderUUID == "" {
		return apperrors.ErrInvalidRequest
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	orderData, exists := r.storage[order.OrderUUID]
	if !exists {
		return apperrors.ErrNotFound
	}

	if orderData.Order.Status == string(model.OrderStatusCANCELLED) {
		return apperrors.ErrForbidden
	}

	if orderData.Order.Status == string(model.OrderStatusPAID) &&
		order.Status != string(model.OrderStatusCANCELLED) {
		return apperrors.ErrForbidden
	}

	// Update order keeping parts
	orderData.Order = order
	r.storage[order.OrderUUID] = orderData

	return nil
}

func (r *inMemoryRepo) DeleteOrder(ctx context.Context, uuid string) error {
	if uuid == "" {
		return apperrors.ErrInvalidRequest
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	orderData, exists := r.storage[uuid]
	if !exists {
		return apperrors.ErrNotFound
	}

	if orderData.Order.Status == string(model.OrderStatusCANCELLED) {
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

	orderData, exists := r.storage[uuid]
	if !exists {
		return apperrors.ErrNotFound
	}

	if orderData.Order.Status == string(model.OrderStatusCANCELLED) {
		return apperrors.ErrForbidden
	}

	if orderData.Order.Status == string(model.OrderStatusPAID) {
		return apperrors.ErrForbidden
	}

	orderData.Order.Status = string(model.OrderStatusCANCELLED)
	r.storage[uuid] = orderData

	return nil
}
