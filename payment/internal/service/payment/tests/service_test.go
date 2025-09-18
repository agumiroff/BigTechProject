package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/agumiroff/BigTechProject/payment/v1/internal/model"
	repoErrors "github.com/agumiroff/BigTechProject/payment/v1/internal/repository/payment"
	paymentService "github.com/agumiroff/BigTechProject/payment/v1/internal/service/payment"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) PayOrder(ctx context.Context, p *model.Payment) (string, error) {
	args := m.Called(ctx, p)
	return args.String(0), args.Error(1)
}

func (m *mockRepo) GetPayment(ctx context.Context, uuid string) (*model.Payment, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Payment), args.Error(1)
}

func TestService_PayOrder_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mockRepo)
	svc := paymentService.NewService(mockRepo)

	ctx := context.Background()
	payment := &model.Payment{
		UUID:          "user-123",
		OrderUUID:     "order-123",
		PaymentMethod: model.PaymentMethodCard,
	}
	expectedTxID := "tx-123"

	mockRepo.On("PayOrder", ctx, payment).Return(expectedTxID, nil)

	// Act
	txID, err := svc.PayOrder(ctx, payment)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedTxID, txID)
	mockRepo.AssertExpectations(t)
}

func TestService_PayOrder_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(mockRepo)
	svc := paymentService.NewService(mockRepo)

	ctx := context.Background()
	payment := &model.Payment{
		UUID:          "user-123",
		OrderUUID:     "order-123",
		PaymentMethod: model.PaymentMethodCard,
	}

	expectedErr := repoErrors.ErrPaymentRequired
	mockRepo.On("PayOrder", ctx, payment).Return("", expectedErr)

	// Act
	txID, err := svc.PayOrder(ctx, payment)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "", txID)
	assert.True(t, errors.Is(err, expectedErr))
	mockRepo.AssertExpectations(t)
}

func TestService_PayOrder_ValidationErrors(t *testing.T) {
	testCases := []struct {
		name     string
		payment  *model.Payment
		mockFunc func(*mockRepo)
		expErr   error
	}{
		{
			name:     "nil payment",
			payment:  nil,
			mockFunc: func(r *mockRepo) {},
			expErr:   repoErrors.ErrPaymentRequired,
		},

		{
			name: "invalid payment method",
			payment: &model.Payment{
				UUID:          "user-123",
				OrderUUID:     "order-123",
				PaymentMethod: "",
			},
			mockFunc: func(r *mockRepo) {},
			expErr:   repoErrors.ErrPaymentMethodInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(mockRepo)
			tc.mockFunc(mockRepo)
			svc := paymentService.NewService(mockRepo)

			// Act
			txID, err := svc.PayOrder(context.Background(), tc.payment)

			// Assert
			assert.Error(t, err)
			assert.Equal(t, "", txID)
			assert.True(t, errors.Is(err, tc.expErr))
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_PayOrder_ContextCanceled(t *testing.T) {
	// Arrange
	mockRepo := new(mockRepo)
	svc := paymentService.NewService(mockRepo)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel context immediately

	payment := &model.Payment{
		UUID:          "user-123",
		OrderUUID:     "order-123",
		PaymentMethod: model.PaymentMethodCard,
	}

	mockRepo.On("PayOrder", ctx, payment).Return("", context.Canceled)

	// Act
	txID, err := svc.PayOrder(ctx, payment)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "", txID)
	assert.True(t, errors.Is(err, context.Canceled))
	mockRepo.AssertExpectations(t)
}

func TestService_PayOrder_ContextDeadlineExceeded(t *testing.T) {
	// Arrange
	mockRepo := new(mockRepo)
	svc := paymentService.NewService(mockRepo)

	ctx, cancel := context.WithTimeout(context.Background(), 0) // Immediate timeout
	defer cancel()

	payment := &model.Payment{
		UUID:          "user-123",
		OrderUUID:     "order-123",
		PaymentMethod: model.PaymentMethodCard,
	}

	mockRepo.On("PayOrder", ctx, payment).Return("", context.DeadlineExceeded)

	// Act
	txID, err := svc.PayOrder(ctx, payment)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "", txID)
	assert.True(t, errors.Is(err, context.DeadlineExceeded))
	mockRepo.AssertExpectations(t)
}
