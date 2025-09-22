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

func TestService_PayOrder_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mockRepo)
	svc := paymentService.NewService(mockRepo)

	ctx := context.Background()
	payment := &model.Payment{
		UserUuid:      "user-123",
		OrderUuid:     "order-123",
		PaymentMethod: model.CARD,
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
		UserUuid:      "user-123",
		OrderUuid:     "order-123",
		PaymentMethod: model.CARD,
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
			name:    "nil payment",
			payment: nil,
			mockFunc: func(r *mockRepo) {
				r.On("PayOrder", mock.Anything, (*model.Payment)(nil)).
					Return("", repoErrors.ErrPaymentRequired)
			},
			expErr: repoErrors.ErrPaymentRequired,
		},
		{
			name: "empty user uuid",
			payment: &model.Payment{
				OrderUuid:     "order-123",
				PaymentMethod: model.CARD,
			},
			mockFunc: func(r *mockRepo) {
				r.On("PayOrder", mock.Anything, mock.MatchedBy(func(p *model.Payment) bool {
					return p.OrderUuid == "order-123" && p.UserUuid == "" && p.PaymentMethod == model.CARD
				})).Return("", repoErrors.ErrUserUUIDRequired)
			},
			expErr: repoErrors.ErrUserUUIDRequired,
		},
		{
			name: "empty order uuid",
			payment: &model.Payment{
				UserUuid:      "user-123",
				PaymentMethod: model.CARD,
			},
			mockFunc: func(r *mockRepo) {
				r.On("PayOrder", mock.Anything, mock.MatchedBy(func(p *model.Payment) bool {
					return p.UserUuid == "user-123" && p.OrderUuid == "" && p.PaymentMethod == model.CARD
				})).Return("", repoErrors.ErrOrderUUIDRequired)
			},
			expErr: repoErrors.ErrOrderUUIDRequired,
		},
		{
			name: "invalid payment method",
			payment: &model.Payment{
				UserUuid:      "user-123",
				OrderUuid:     "order-123",
				PaymentMethod: model.CategoryUnspecified,
			},
			mockFunc: func(r *mockRepo) {
				r.On("PayOrder", mock.Anything, mock.MatchedBy(func(p *model.Payment) bool {
					return p.UserUuid == "user-123" && p.OrderUuid == "order-123" && p.PaymentMethod == model.CategoryUnspecified
				})).Return("", repoErrors.ErrPaymentMethodInvalid)
			},
			expErr: repoErrors.ErrPaymentMethodInvalid,
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
		UserUuid:      "user-123",
		OrderUuid:     "order-123",
		PaymentMethod: model.CARD,
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
		UserUuid:      "user-123",
		OrderUuid:     "order-123",
		PaymentMethod: model.CARD,
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
