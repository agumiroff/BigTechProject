package inmemory

import (
	"context"
	"fmt"

	"github.com/agumiroff/BigTechProject/shared/apperrors"
	"sync"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	repomodel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
)

type InMemoryOrderRepository struct {
	mu     sync.RWMutex
	orders map[string]*repomodel.Order
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[string]*repomodel.Order),
	}
}

func (r *InMemoryOrderRepository) CreateOrder(ctx context.Context, order *model.Order) (*model.CreateOrderResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if order == nil {
		return nil, apperrors.ErrInvalidRequest
	}

	if order.OrderUUID == "" {
		return nil, apperrors.ErrInvalidRequest
	}

	if order.UserUUID == "" {
		return nil, apperrors.ErrInvalidRequest
	}

	if len(order.PartUUIDs) == 0 {
		return nil, apperrors.ErrInvalidRequest
	}

	if _, exists := r.orders[order.OrderUUID]; exists {
		return nil, apperrors.ErrAlreadyExists
	}

	paymentMethod := repomodel.PaymentMethodUNKNOWN
	if order.PaymentMethod != "" {
		paymentMethod = repomodel.PaymentMethod(order.PaymentMethod)
	}

	repoOrder := &repomodel.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          repomodel.OrderStatus(order.Status),
	}

	r.orders[order.OrderUUID] = repoOrder

	return &model.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}

func (r *InMemoryOrderRepository) Get(ctx context.Context, uuid string) (*repomodel.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if uuid == "" {
		return nil, apperrors.ErrInvalidRequest
	}

	order, exists := r.orders[uuid]
	if !exists {
		return nil, apperrors.ErrNotFound
	}

	return order, nil
}

func (r *InMemoryOrderRepository) UpdateOrder(ctx context.Context, order *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[order.OrderUUID]; !exists {
		return fmt.Errorf("order with UUID %s not found", order.OrderUUID)
	}

	paymentMethod := repomodel.PaymentMethodUNKNOWN
	if order.PaymentMethod != "" {
		paymentMethod = repomodel.PaymentMethod(order.PaymentMethod)
	}

	r.orders[order.OrderUUID] = &repomodel.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          repomodel.OrderStatus(order.Status),
	}

	return nil
}

func (r *InMemoryOrderRepository) DeleteOrder(ctx context.Context, uuid string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if uuid == "" {
		return apperrors.ErrInvalidRequest
	}

	order, exists := r.orders[uuid]
	if !exists {
		return apperrors.ErrNotFound
	}

	if order.Status == repomodel.OrderStatusCANCELLED {
		return apperrors.ErrForbidden
	}

	delete(r.orders, uuid)
	return nil
}

func (r *InMemoryOrderRepository) CancelOrder(ctx context.Context, uuid string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, exists := r.orders[uuid]
	if !exists {
		return fmt.Errorf("order with UUID %s not found", uuid)
	}

	order.Status = repomodel.OrderStatusCANCELLED
	return nil
}
