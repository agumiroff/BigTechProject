package service

import (
	"context"
	"log"
	"sync"

	"github.com/brianvoe/gofakeit/v6"

	PayServiceV1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

type PayService struct {
	PayServiceV1.UnimplementedPaymentServiceServer
	mu sync.RWMutex

	storage map[string]*PayServiceV1.Payment
}

func NewService() (res *PayService) {
	service := &PayService{
		storage: make(map[string]*PayServiceV1.Payment),
	}

	return service
}

func (s *PayService) PayOrder(ctx context.Context, req *PayServiceV1.PayOrderRequest) (res *PayServiceV1.PayOrderResponse, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id := gofakeit.UUID()

	s.storage[id] = req.Payment

	log.Printf("Payment successfully created: %s\n", id)

	response := &PayServiceV1.PayOrderResponse{
		TransactionUuid: id,
	}

	return response, nil
}
