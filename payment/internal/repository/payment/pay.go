package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/agumiroff/BigTechProject/payment/v1/internal/model"
	"github.com/agumiroff/BigTechProject/payment/v1/internal/repository/converter"
)

func (r *repository) PayOrder(ctx context.Context, p *model.Payment) (string, error) {
	if p == nil {
		return "", ErrPaymentRequired
	}

	if p.UUID == "" {
		return "", ErrUserUUIDRequired
	}

	if p.OrderUUID == "" {
		return "", ErrOrderUUIDRequired
	}

	if p.PaymentMethod == "" {
		return "", ErrPaymentMethodInvalid
	}

	// Convert domain model to repository model
	repoPayment := converter.ModelToRepo(p)
	txid := gofakeit.UUID()
	repoPayment.UUID = txid
	repoPayment.CreatedAt = time.Now()
	repoPayment.UpdatedAt = time.Now()

	// Insert payment
	_, err := r.collection.InsertOne(ctx, repoPayment)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", fmt.Errorf("payment for order %s already exists", p.OrderUUID)
		}
		return "", fmt.Errorf("failed to create payment: %w", err)
	}

	return txid, nil
}
