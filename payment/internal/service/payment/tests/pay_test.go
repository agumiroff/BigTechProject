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

type mockPaymentRepository struct {
	mock.Mock
}

func (m *mockPaymentRepository) PayOrder(ctx context.Context, p *model.Payment) (string, error) {
	args := m.Called(ctx, p)
	return args.String(0), args.Error(1)
}

func (m *mockPaymentRepository) GetPayment(ctx context.Context, uuid string) (*model.Payment, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Payment), args.Error(1)
}

func TestPayOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		mockRepo := new(mockPaymentRepository)
		svc := paymentService.NewService(mockRepo)

		payment := &model.Payment{
			UUID:          "user-123",
			OrderUUID:     "order-123",
			PaymentMethod: model.PaymentMethodCard,
		}
		expectedTxID := "tx-123"

		mockRepo.On("PayOrder", mock.Anything, payment).Return(expectedTxID, nil)

		// Act
		txID, err := svc.PayOrder(context.Background(), payment)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedTxID, txID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation errors", func(t *testing.T) {
		testCases := []struct {
			name    string
			payment *model.Payment
			expErr  error
		}{
			{
				name:    "nil payment",
				payment: nil,
				expErr:  repoErrors.ErrPaymentRequired,
			},

			{
				name: "unspecified payment method",
				payment: &model.Payment{
					UUID:      "test-user",
					OrderUUID: "test-order",
				},
				expErr: repoErrors.ErrPaymentMethodInvalid,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Arrange
				mockRepo := new(mockPaymentRepository)

				svc := paymentService.NewService(mockRepo)

				// Act
				txID, err := svc.PayOrder(context.Background(), tc.payment)

				// Assert
				assert.Error(t, err)
				assert.Equal(t, "", txID)
				assert.ErrorIs(t, err, tc.expErr)
				// Validation happens before mock is called
			})
		}
	})

	t.Run("context handling", func(t *testing.T) {
		testCases := []struct {
			name    string
			ctxFunc func() (context.Context, context.CancelFunc)
			expErr  error
		}{
			{
				name: "context canceled",
				ctxFunc: func() (context.Context, context.CancelFunc) {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx, cancel
				},
				expErr: context.Canceled,
			},
			{
				name: "context deadline exceeded",
				ctxFunc: func() (context.Context, context.CancelFunc) {
					return context.WithTimeout(context.Background(), 0)
				},
				expErr: context.DeadlineExceeded,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Arrange
				mockRepo := new(mockPaymentRepository)
				svc := paymentService.NewService(mockRepo)

				ctx, cancel := tc.ctxFunc()
				defer cancel()

				payment := &model.Payment{
					UUID:          "user-123",
					OrderUUID:     "order-123",
					PaymentMethod: model.PaymentMethodCard,
				}

				mockRepo.On("PayOrder", ctx, payment).Return("", tc.expErr)

				// Act
				txID, err := svc.PayOrder(ctx, payment)

				// Assert
				assert.Error(t, err)
				assert.Equal(t, "", txID)
				assert.True(t, errors.Is(err, tc.expErr))
				mockRepo.AssertExpectations(t)
			})
		}
	})
}

func TestGetPayment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		mockRepo := new(mockPaymentRepository)
		svc := paymentService.NewService(mockRepo)

		expectedPayment := &model.Payment{
			UUID:          "user-123",
			OrderUUID:     "order-123",
			PaymentMethod: model.PaymentMethodCard,
		}
		paymentID := "payment-123"

		mockRepo.On("GetPayment", mock.Anything, paymentID).Return(expectedPayment, nil)

		// Act
		payment, err := svc.GetPayment(context.Background(), paymentID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedPayment, payment)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(mockPaymentRepository)
		svc := paymentService.NewService(mockRepo)

		paymentID := "non-existent-payment"
		mockRepo.On("GetPayment", mock.Anything, paymentID).Return(nil, repoErrors.ErrPaymentNotFound)

		// Act
		payment, err := svc.GetPayment(context.Background(), paymentID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.ErrorIs(t, err, repoErrors.ErrPaymentNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("context handling", func(t *testing.T) {
		testCases := []struct {
			name    string
			ctxFunc func() (context.Context, context.CancelFunc)
			expErr  error
		}{
			{
				name: "context canceled",
				ctxFunc: func() (context.Context, context.CancelFunc) {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx, cancel
				},
				expErr: context.Canceled,
			},
			{
				name: "context deadline exceeded",
				ctxFunc: func() (context.Context, context.CancelFunc) {
					return context.WithTimeout(context.Background(), 0)
				},
				expErr: context.DeadlineExceeded,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Arrange
				mockRepo := new(mockPaymentRepository)
				svc := paymentService.NewService(mockRepo)

				ctx, cancel := tc.ctxFunc()
				defer cancel()

				paymentID := "payment-123"
				mockRepo.On("GetPayment", ctx, paymentID).Return(nil, tc.expErr)

				// Act
				payment, err := svc.GetPayment(ctx, paymentID)

				// Assert
				assert.Error(t, err)
				assert.Nil(t, payment)
				assert.True(t, errors.Is(err, tc.expErr))
				mockRepo.AssertExpectations(t)
			})
		}
	})
}
