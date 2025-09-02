package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/agumiroff/BigTechProject/payment/v1/internal/model"
	"github.com/agumiroff/BigTechProject/payment/v1/internal/repository/payment"
)

func newTestPayment() *model.Payment {
	return &model.Payment{
		UserUuid:      "test-user-uuid",
		OrderUuid:     "test-order-uuid",
		PaymentMethod: model.CARD,
	}
}

func TestPayOrder_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := payment.NewRepository()
	testPayment := newTestPayment()

	// Act
	txID, err := repo.PayOrder(ctx, testPayment)

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, txID)

	// Verify payment was stored correctly
	stored, err := repo.GetPayment(ctx, txID)
	require.NoError(t, err)
	require.NotNil(t, stored)

	assert.Equal(t, testPayment.UserUuid, stored.UserUuid)
	assert.Equal(t, testPayment.OrderUuid, stored.OrderUuid)
	assert.Equal(t, testPayment.PaymentMethod, stored.PaymentMethod)
}

func TestPayOrder_ValidationErrors(t *testing.T) {
	testCases := []struct {
		name    string
		payment *model.Payment
		errMsg  string
	}{
		{
			name:    "nil payment",
			payment: nil,
			errMsg:  "payment is required",
		},
		{
			name: "empty user uuid",
			payment: &model.Payment{
				OrderUuid:     "test-order",
				PaymentMethod: model.CARD,
			},
			errMsg: "user uuid is required",
		},
		{
			name: "empty order uuid",
			payment: &model.Payment{
				UserUuid:      "test-user",
				PaymentMethod: model.CARD,
			},
			errMsg: "order uuid is required",
		},
		{
			name: "unspecified payment method",
			payment: &model.Payment{
				UserUuid:      "test-user",
				OrderUuid:     "test-order",
				PaymentMethod: model.CategoryUnspecified,
			},
			errMsg: "payment method is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			repo := payment.NewRepository()

			// Act
			txID, err := repo.PayOrder(ctx, tc.payment)

			// Assert
			require.Error(t, err)
			require.Empty(t, txID)
			assert.Contains(t, err.Error(), tc.errMsg)
		})
	}
}

func TestGetPayment_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := payment.NewRepository()
	testPayment := newTestPayment()

	// Create payment first
	txID, err := repo.PayOrder(ctx, testPayment)
	require.NoError(t, err)
	require.NotEmpty(t, txID)

	// Act
	stored, err := repo.GetPayment(ctx, txID)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, stored)

	assert.Equal(t, testPayment.UserUuid, stored.UserUuid)
	assert.Equal(t, testPayment.OrderUuid, stored.OrderUuid)
	assert.Equal(t, testPayment.PaymentMethod, stored.PaymentMethod)
}

func TestGetPayment_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := payment.NewRepository()
	nonExistentTxID := "non-existent-tx-id"

	// Act
	stored, err := repo.GetPayment(ctx, nonExistentTxID)

	// Assert
	require.Error(t, err)
	require.Nil(t, stored)
	assert.ErrorIs(t, err, payment.ErrPaymentNotFound)
}

func TestGetPayment_EmptyTxID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := payment.NewRepository()

	// Act
	stored, err := repo.GetPayment(ctx, "")

	// Assert
	require.Error(t, err)
	require.Nil(t, stored)
	assert.ErrorIs(t, err, payment.ErrTxIDRequired)
}
