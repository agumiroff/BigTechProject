package order

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/agumiroff/BigTechProject/order/v1/internal/converter"
	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	"github.com/agumiroff/BigTechProject/platform/pkg/grpc/logger"
	inventoryv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
)

func (s *service) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.CreateOrderResponse, error) {
	// Validate request
	if err := validateCreateOrderRequest(req); err != nil {
		logger.Warn(ctx, "Invalid create order request", zap.Error(err))
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Get parts from inventory service
	list, err := s.ExRepo.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartsFilter{
			Uuids: req.PartUUIDs,
		},
	})
	if err != nil {
		logger.Error(ctx, "Failed to list parts from inventory", zap.Error(err))
		return nil, fmt.Errorf("failed to get parts: %w", err)
	}
	logger.Debug(ctx, "list of parts successfully loaded", zap.Int("parts_count", len(list.Parts)))

	// Validate all parts exist
	if err = validatePartsExist(req.PartUUIDs, list.Parts); err != nil {
		logger.Warn(ctx, "parts validation failed", zap.Error(err))
		return nil, err
	}

	// Calculate total price and extract part UUIDs
	totalPrice, partUUIDs := calculateOrderDetails(list.Parts)

	// Create order
	order := &model.Order{
		UserUUID:   req.UserUUID,
		PartUUIDs:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPENDINGPAYMENT,
	}

	// Save order
	orderRow := converter.ToRepoOrder(order)
	id, err := s.Repo.CreateOrder(ctx, orderRow, order.PartUUIDs)
	if err != nil {
		logger.Error(ctx, "Failed to create order in repository", zap.Error(err))
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	logger.Info(ctx, "Order created successfully", zap.String("order_uuid", order.OrderUUID), zap.Float64("total_price", order.TotalPrice))

	return &model.CreateOrderResponse{
		OrderUUID:  id,
		TotalPrice: totalPrice,
	}, nil
}

func validateCreateOrderRequest(req *model.CreateOrderRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if req.UserUUID == "" {
		return fmt.Errorf("user uuid is required")
	}
	if len(req.PartUUIDs) == 0 {
		return fmt.Errorf("at least one part uuid is required")
	}
	return nil
}

func validatePartsExist(requestedUUIDs []string, parts []*inventoryv1.Part) error {
	// Create a map of returned part UUIDs for O(1) lookup
	partMap := make(map[string]bool)
	for _, part := range parts {
		partMap[part.GetUuid()] = true
	}

	// Check each requested UUID exists in returned parts
	for _, uuid := range requestedUUIDs {
		if !partMap[uuid] {
			return fmt.Errorf("parts not found: %s", uuid)
		}
	}

	return nil
}

func calculateOrderDetails(parts []*inventoryv1.Part) (totalPrice float64, partUUIDs []string) {
	for _, part := range parts {
		totalPrice += part.GetPrice()
		partUUIDs = append(partUUIDs, part.GetUuid())
	}
	return totalPrice, partUUIDs
}
