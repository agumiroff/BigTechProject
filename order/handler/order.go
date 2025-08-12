package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"

	OrderV1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
	InvV1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
	PayV1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

type OrderStorage struct {
	mu sync.RWMutex

	storage map[string]*OrderV1.Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		storage: make(map[string]*OrderV1.Order),
	}
}

func mapPaymentMethod(method OrderV1.PaymentMethod) PayV1.PaymentMethod {
	switch method {
	case OrderV1.PaymentMethodCARD:
		return PayV1.PaymentMethod_PAYMENT_METHOD_CARD
	case OrderV1.PaymentMethodSBP:
		return PayV1.PaymentMethod_PAYMENT_METHOD_SBP
	case OrderV1.PaymentMethodCREDITCARD:
		return PayV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case OrderV1.PaymentMethodINVESTORMONEY:
		return PayV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return PayV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

func (s *OrderStorage) CreateOrder(orderUuid string, order *OrderV1.Order) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if order, exist := s.storage[orderUuid]; exist {
		log.Printf("order with %s is already in storage: \n%v", orderUuid, order)
		return ""
	}

	s.storage[orderUuid] = order
	return order.GetOrderUUID()
}

func (s *OrderStorage) GetOrder(orderUuid string) *OrderV1.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, exist := s.storage[orderUuid]
	if !exist {
		log.Printf("Order with %s not found", orderUuid)
		return nil
	}

	return order
}

func (s *OrderStorage) DeleteOrder(orderUuid string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	delete(s.storage, orderUuid)
	log.Printf("value deleted")
}

func (s *OrderStorage) UpdateOrder(orderUuid string, order *OrderV1.Order) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.storage[orderUuid] = order
	log.Printf("value updated")
}

type OrderHandler struct {
	storage *OrderStorage

	invClient InvV1.InventoryServiceClient
	payClient PayV1.PaymentServiceClient
}

func NewOrderHandler(invClient InvV1.InventoryServiceClient, payClient PayV1.PaymentServiceClient) *OrderHandler {
	storage := NewOrderStorage()
	return &OrderHandler{
		storage:   storage,
		invClient: invClient,
		payClient: payClient,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *OrderV1.CreateOrderRequest) (res OrderV1.CreateOrderRes, err error) {
	sum := 0.0
	orderUuid := gofakeit.UUID()
	// Check all parts in invService

	uuids := req.GetPartUuids()

	allPartsResp, err := h.invClient.ListParts(ctx, &InvV1.ListPartsRequest{
		Filter: &InvV1.PartsFilter{
			Uuids: uuids,
		},
	})
	if err != nil {
		log.Printf("failed to list parts: %s", err)
		return nil, err
	}

	allParts := allPartsResp.Parts
	var allPartsUuids []string
	for _, part := range allParts {
		allPartsUuids = append(allPartsUuids, part.GetUuid())
	}

	// Check all parts are exist, if no - return nil
	log.Printf("Checking all parts are exist")
	if len(allPartsResp.GetParts()) != len(req.GetPartUuids()) {
		log.Printf("Some parts not found")
		return nil, fmt.Errorf("some parts not found")
	}

	// Calculating total price
	for _, p := range allParts {
		sum += p.GetPrice()
	}

	// Generate order UUID
	uuid := gofakeit.UUID()

	// Save order with status PENDING PAYMENT
	order := &OrderV1.Order{
		OrderUUID:  orderUuid,
		UserUUID:   req.GetUserUUID(),
		PartUuids:  allPartsUuids,
		TotalPrice: sum,
		Status:     OrderV1.OrderStatusPENDINGPAYMENT,
	}
	h.storage.CreateOrder(order.GetOrderUUID(), order)

	log.Printf("Order created uuid: %s\n, sum: %v\n", uuid, sum)

	return &OrderV1.CreateOrderResponse{
		OrderUUID:  orderUuid,
		TotalPrice: sum,
	}, nil
}

func (h *OrderHandler) CancelOrderByUuid(ctx context.Context, params OrderV1.CancelOrderByUuidParams) (OrderV1.CancelOrderByUuidRes, error) {
	orderUUID := params.OrderUUID
	h.storage.DeleteOrder(orderUUID.String())

	log.Printf("order UUID %d", orderUUID)

	return &OrderV1.CancelOrderResponse{
		OrderUUID: orderUUID,
	}, nil
}

func (h *OrderHandler) GetOrderByUuid(ctx context.Context, params OrderV1.GetOrderByUuidParams) (OrderV1.GetOrderByUuidRes, error) {
	order := h.storage.GetOrder(params.OrderUUID.String())

	log.Printf("order found %+v", order)
	return &OrderV1.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   order.PaymentMethod,
		Status:          order.Status,
	}, nil
}

func (h *OrderHandler) PayOrder(ctx context.Context, req *OrderV1.PayOrderRequest, params OrderV1.PayOrderParams) (OrderV1.PayOrderRes, error) {
	order := h.storage.GetOrder(params.OrderUUID.String())

	_, err := h.payClient.PayOrder(ctx,
		&PayV1.PayOrderRequest{
			Payment: &PayV1.Payment{
				OrderUuid:     order.GetOrderUUID(),
				UserUuid:      order.GetUserUUID(),
				PaymentMethod: mapPaymentMethod(req.GetPaymentMethod()),
			},
		},
	)
	if err != nil {
		log.Printf("Payment failed %s", err)
		return nil, err
	}

	order.Status = "PAID"

	transUuid, err := uuid.Parse(gofakeit.UUID())
	if err != nil {
		log.Printf("Invalid UUID %d", err)
	}

	return &OrderV1.PayOrderResponse{
		TransactionUUID: transUuid,
	}, nil
}

func (h *OrderHandler) NewError(ctx context.Context, err error) *OrderV1.GenericErrorStatusCode {
	return &OrderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: OrderV1.GenericError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}
}
