package pay_test

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
			UserUuid:      "user-123",
			OrderUuid:     "order-123",
			PaymentMethod: model.CARD,
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
			name        string
			payment     *model.Payment
			setupMock   func(*mockPaymentRepository)
			expectedErr error
		}{
			{
				name:    "nil payment",
				payment: nil,
				setupMock: func(m *mockPaymentRepository) {
					m.On("PayOrder", mock.Anything, (*model.Payment)(nil)).
						Return("", repoErrors.ErrPaymentRequired)
				},
				expectedErr: repoErrors.ErrPaymentRequired,
			},
			{
				name: "empty user uuid",
				payment: &model.Payment{
					OrderUuid:     "order-123",
					PaymentMethod: model.CARD,
				},
				setupMock: func(m *mockPaymentRepository) {
					m.On("PayOrder", mock.Anything, mock.MatchedBy(func(p *model.Payment) bool {
						return p.OrderUuid == "order-123" && p.UserUuid == ""
					})).Return("", repoErrors.ErrUserUUIDRequired)
				},
				expectedErr: repoErrors.ErrUserUUIDRequired,
			},
			{
				name: "empty order uuid",
				payment: &model.Payment{
					UserUuid:      "user-123",
					PaymentMethod: model.CARD,
				},
				setupMock: func(m *mockPaymentRepository) {
					m.On("PayOrder", mock.Anything, mock.MatchedBy(func(p *model.Payment) bool {
						return p.UserUuid == "user-123" && p.OrderUuid == ""
					})).Return("", repoErrors.ErrOrderUUIDRequired)
				},
				expectedErr: repoErrors.ErrOrderUUIDRequired,
			},
			{
				name: "invalid payment method",
				payment: &model.Payment{
					UserUuid:      "user-123",
					OrderUuid:     "order-123",
					PaymentMethod: model.CategoryUnspecified,
				},
				setupMock: func(m *mockPaymentRepository) {
					m.On("PayOrder", mock.Anything, mock.MatchedBy(func(p *model.Payment) bool {
						return p.UserUuid == "user-123" && p.OrderUuid == "order-123" && p.PaymentMethod == model.CategoryUnspecified
					})).Return("", repoErrors.ErrPaymentMethodInvalid)
				},
				expectedErr: repoErrors.ErrPaymentMethodInvalid,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Arrange
				mockRepo := new(mockPaymentRepository)
				tc.setupMock(mockRepo)
				svc := paymentService.NewService(mockRepo)

				// Act
				txID, err := svc.PayOrder(context.Background(), tc.payment)

				// Assert
				assert.Error(t, err)
				assert.Equal(t, "", txID)
				assert.True(t, errors.Is(err, tc.expectedErr))
				mockRepo.AssertExpectations(t)
			})
		}
	})

	t.Run("context handling", func(t *testing.T) {
		testCases := []struct {
			name        string
			setupCtx    func() (context.Context, context.CancelFunc)
			expectedErr error
		}{
			{
				name: "context canceled",
				setupCtx: func() (context.Context, context.CancelFunc) {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx, cancel
				},
				expectedErr: context.Canceled,
			},
			{
				name: "context deadline exceeded",
				setupCtx: func() (context.Context, context.CancelFunc) {
					return context.WithTimeout(context.Background(), 0)
				},
				expectedErr: context.DeadlineExceeded,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Arrange
				mockRepo := new(mockPaymentRepository)
				svc := paymentService.NewService(mockRepo)

				ctx, cancel := tc.setupCtx()
				defer cancel()

				payment := &model.Payment{
					UserUuid:      "user-123",
					OrderUuid:     "order-123",
					PaymentMethod: model.CARD,
				}

				mockRepo.On("PayOrder", ctx, payment).Return("", tc.expectedErr)

				// Act
				txID, err := svc.PayOrder(ctx, payment)

				// Assert
				assert.Error(t, err)
				assert.Equal(t, "", txID)
				assert.True(t, errors.Is(err, tc.expectedErr))
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
			UserUuid:      "user-123",
			OrderUuid:     "order-123",
			PaymentMethod: model.CARD,
		}
		txID := "tx-123"

		mockRepo.On("GetPayment", mock.Anything, txID).Return(expectedPayment, nil)

		// Act
		payment, err := svc.GetPayment(context.Background(), txID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedPayment, payment)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation errors", func(t *testing.T) {
		testCases := []struct {
			name        string
			txID        string
			setupMock   func(*mockPaymentRepository)
			expectedErr error
		}{
			{
				name: "empty transaction id",
				txID: "",
				setupMock: func(m *mockPaymentRepository) {
					m.On("GetPayment", mock.Anything, "").
						Return(nil, repoErrors.ErrTxIDRequired)
				},
				expectedErr: repoErrors.ErrTxIDRequired,
			},
			{
				name: "payment not found",
				txID: "non-existent-tx",
				setupMock: func(m *mockPaymentRepository) {
					m.On("GetPayment", mock.Anything, "non-existent-tx").
						Return(nil, repoErrors.ErrPaymentNotFound)
				},
				expectedErr: repoErrors.ErrPaymentNotFound,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Arrange
				mockRepo := new(mockPaymentRepository)
				tc.setupMock(mockRepo)
				svc := paymentService.NewService(mockRepo)

				// Act
				payment, err := svc.GetPayment(context.Background(), tc.txID)

				// Assert
				assert.Error(t, err)
				assert.Nil(t, payment)
				assert.True(t, errors.Is(err, tc.expectedErr))
				mockRepo.AssertExpectations(t)
			})
		}
	})

	t.Run("context handling", func(t *testing.T) {
		testCases := []struct {
			name        string
			setupCtx    func() (context.Context, context.CancelFunc)
			expectedErr error
		}{
			{
				name: "context canceled",
				setupCtx: func() (context.Context, context.CancelFunc) {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx, cancel
				},
				expectedErr: context.Canceled,
			},
			{
				name: "context deadline exceeded",
				setupCtx: func() (context.Context, context.CancelFunc) {
					return context.WithTimeout(context.Background(), 0)
				},
				expectedErr: context.DeadlineExceeded,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Arrange
				mockRepo := new(mockPaymentRepository)
				svc := paymentService.NewService(mockRepo)

				ctx, cancel := tc.setupCtx()
				defer cancel()

				txID := "tx-123"
				mockRepo.On("GetPayment", ctx, txID).Return(nil, tc.expectedErr)

				// Act
				payment, err := svc.GetPayment(ctx, txID)

				// Assert
				assert.Error(t, err)
				assert.Nil(t, payment)
				assert.True(t, errors.Is(err, tc.expectedErr))
				mockRepo.AssertExpectations(t)
			})
		}
	})
}
