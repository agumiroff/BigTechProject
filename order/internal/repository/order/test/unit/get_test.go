package unit

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	repomodel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func TestGetOrder_Success(t *testing.T) {
	// Arrange
	repo := NewInmemoryRepo()
	ctx := context.Background()

	orderUUID := "test-order-uuid"
	userUUID := "test-user-uuid"
	parts := []string{"part1", "part2"}
	totalPrice := 100.0

	order := &repomodel.OrderRow{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		TotalPrice: totalPrice,
		Status:     string(model.OrderStatusPENDINGPAYMENT),
	}

	// Create the order
	_, err := repo.CreateOrder(ctx, order, parts)
	require.NoError(t, err)

	// Act
	gotOrder, gotParts, err := repo.GetOrder(ctx, orderUUID)
	_ = gotParts // Prevent unused variable warning

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, gotOrder)
	assert.Equal(t, orderUUID, gotOrder.OrderUUID)
	assert.Equal(t, userUUID, gotOrder.UserUUID)
	assert.Equal(t, totalPrice, gotOrder.TotalPrice)
	assert.Equal(t, string(model.OrderStatusPENDINGPAYMENT), gotOrder.Status)
	assert.Equal(t, parts, gotParts)
}

func TestGetOrder_NotFound(t *testing.T) {
	// Arrange
	repo := NewInmemoryRepo()
	ctx := context.Background()

	// Act
	_, _, err := repo.GetOrder(ctx, "nonexistent-uuid")

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrNotFound)
}

func TestGetOrder_EmptyUUID(t *testing.T) {
	// Arrange
	repo := NewInmemoryRepo()
	ctx := context.Background()

	// Act
	_, _, err := repo.GetOrder(ctx, "")

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrInvalidRequest)
}

func TestGetOrder_WithTransactionAndPaymentInfo(t *testing.T) {
	// Arrange
	repo := NewInmemoryRepo()
	ctx := context.Background()

	orderUUID := "test-order-uuid"
	userUUID := "test-user-uuid"
	parts := []string{"part1"}
	totalPrice := 100.0
	transactionUUID := "test-transaction-uuid"
	paymentMethod := string(model.PaymentMethodCARD)

	order := &repomodel.OrderRow{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		TotalPrice:      totalPrice,
		Status:          string(model.OrderStatusPAID),
		TransactionUUID: sql.NullString{String: transactionUUID, Valid: true},
		PaymentMethod:   sql.NullString{String: paymentMethod, Valid: true},
	}

	// Create the order
	_, err := repo.CreateOrder(ctx, order, parts)
	require.NoError(t, err)

	// Act
	gotOrder, gotParts, err := repo.GetOrder(ctx, orderUUID)
	_ = gotParts // Prevent unused variable warning

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, gotOrder)
	assert.Equal(t, orderUUID, gotOrder.OrderUUID)
	assert.Equal(t, string(model.OrderStatusPAID), gotOrder.Status)
	assert.True(t, gotOrder.TransactionUUID.Valid)
	assert.Equal(t, transactionUUID, gotOrder.TransactionUUID.String)
	assert.True(t, gotOrder.PaymentMethod.Valid)
	assert.Equal(t, paymentMethod, gotOrder.PaymentMethod.String)
}
