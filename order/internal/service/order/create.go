package order

import (
	"context"
	"fmt"
	"log"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	inventoryv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
)

func (s *service) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (res *model.CreateOrderResponse, err error) {
	sum := 0.0
	// Check all parts in invService
	uuids := req.PartUUIDs

	list, err := s.ExRepo.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartsFilter{
			Uuids: uuids,
		},
	})
	if err != nil {
		log.Printf("failed to list parts: %s", err)
		return nil, err
	}

	var stockParts []*inventoryv1.Part
	var stockPartUuids []string
	for _, part := range list.Parts {
		stockParts = append(stockParts, part)
		stockPartUuids = append(stockPartUuids, part.Uuid)
	}

	// Check all parts are exist, if no - return nil
	log.Printf("Checking all parts are exist")
	if len(stockParts) != len(uuids) {
		log.Printf("Some parts not found")
		return nil, fmt.Errorf("some parts not found")
	}

	// Calculating total price
	for _, p := range stockParts {
		sum += p.GetPrice()
	}

	// Generate order UUID
	uuid := gofakeit.UUID()

	// Save order with status PENDING PAYMENT
	order := &model.Order{
		UserUUID:   req.UserUUID,
		OrderUUID:  uuid,
		PartUUIDs:  stockPartUuids,
		TotalPrice: sum,
		Status:     model.OrderStatusPENDINGPAYMENT,
	}

	response, err := s.Repo.CreateOrder(ctx, order)
	if err != nil {
		log.Printf("Failed to create order: %v", err)
		return nil, err
	}

	log.Printf("Order created uuid: %s\n, sum: %v\n", uuid, sum)

	return &model.CreateOrderResponse{
		OrderUUID:  response.OrderUUID,
		TotalPrice: response.TotalPrice,
	}, nil
}
