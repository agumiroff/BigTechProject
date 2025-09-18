package payment

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/agumiroff/BigTechProject/payment/v1/internal/model"
	"github.com/agumiroff/BigTechProject/payment/v1/internal/repository/converter"
	"github.com/brianvoe/gofakeit/v6"
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

	log.Printf("payment successfully created")
	return txid, nil
}
