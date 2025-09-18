package payment

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/agumiroff/BigTechProject/payment/v1/internal/model"
	"github.com/agumiroff/BigTechProject/payment/v1/internal/repository/converter"
	repomodel "github.com/agumiroff/BigTechProject/payment/v1/internal/repository/model"
)

func (r *repository) GetPayment(ctx context.Context, uuid string) (*model.Payment, error) {
	if uuid == "" {
		return nil, ErrTxIDRequired
	}

	var payment repomodel.Payment
	err := r.collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&payment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrPaymentNotFound
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return converter.RepoToModel(&payment), nil
}
