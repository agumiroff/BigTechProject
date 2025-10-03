package payment_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/agumiroff/BigTechProject/payment/v1/internal/model"
	"github.com/agumiroff/BigTechProject/payment/v1/internal/repository/payment"
)

func setupTestDB(t *testing.T) *mongo.Database {
	ctx := context.Background()
	credentials := options.Credential{
		Username: "payment",
		Password: "payment",
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27018").SetAuth(credentials))
	require.NoError(t, err)
	db := client.Database("test_payments")

	// Clean up before each test
	err = db.Collection("payments").Drop(ctx)
	require.NoError(t, err)

	// Create indexes
	indexModel := mongo.IndexModel{
		Keys: map[string]interface{}{
			"order_uuid": 1,
		},
		Options: options.Index().SetUnique(true),
	}
	_, err = db.Collection("payments").Indexes().CreateOne(ctx, indexModel)
	require.NoError(t, err)

	return db
}

func TestPayOrder_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	db := setupTestDB(t)
	repo := payment.NewRepository(db)
	testPayment := &model.Payment{
		UUID:          "user-123",
		OrderUUID:     "order-123",
		PaymentMethod: model.PaymentMethodCard,
	}

	// Act
	txID, err := repo.PayOrder(ctx, testPayment)

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, txID)
}

func TestPayOrder_ValidationErrors(t *testing.T) {
	testCases := []struct {
		name    string
		payment *model.Payment
		expErr  error
	}{
		{
			name:    "nil payment",
			payment: nil,
			expErr:  payment.ErrPaymentRequired,
		},
		{
			name: "empty user uuid",
			payment: &model.Payment{
				OrderUUID:     "order-123",
				PaymentMethod: model.PaymentMethodCard,
			},
			expErr: payment.ErrUserUUIDRequired,
		},
		{
			name: "empty order uuid",
			payment: &model.Payment{
				UUID:          "user-123",
				PaymentMethod: model.PaymentMethodCard,
			},
			expErr: payment.ErrOrderUUIDRequired,
		},
		{
			name: "invalid payment method",
			payment: &model.Payment{
				UUID:          "user-123",
				OrderUUID:     "order-123",
				PaymentMethod: "",
			},
			expErr: payment.ErrPaymentMethodInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			db := setupTestDB(t)
			repo := payment.NewRepository(db)

			// Act
			txID, err := repo.PayOrder(ctx, tc.payment)

			// Assert
			assert.Error(t, err)
			assert.Equal(t, "", txID)
			assert.ErrorIs(t, err, tc.expErr)
		})
	}
}

func TestPayOrder_DuplicateOrder(t *testing.T) {
	// Arrange
	ctx := context.Background()
	db := setupTestDB(t)
	repo := payment.NewRepository(db)

	payment := &model.Payment{
		UUID:          "user-123",
		OrderUUID:     "order-123",
		PaymentMethod: model.PaymentMethodCard,
	}

	// First payment should succeed
	_, err := repo.PayOrder(ctx, payment)
	require.NoError(t, err)

	// Act - Try to create duplicate payment
	txID, err := repo.PayOrder(ctx, payment)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "", txID)
	assert.Contains(t, err.Error(), "payment for order order-123 already exists")
}
