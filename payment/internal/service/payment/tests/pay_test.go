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

func TestPayOrder_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mockRepo)
	svc := paymentService.NewService(mockRepo)

	ctx := context.Background()
	expectedPayment := &model.Payment{
		UserUuid:      "user-123",
		OrderUuid:     "order-123",
		PaymentMethod: model.CARD,
	}
	expectedTxID := "tx-123"

	mockRepo.On("PayOrder", ctx, expectedPayment).
		Return(expectedTxID, nil)

	// Act
	txID, err := svc.PayOrder(ctx, expectedPayment)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedTxID, txID)
	mockRepo.AssertExpectations(t)
}

func TestPayOrder_RepositoryError(t *testing.T) {
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
	mockRepo.On("PayOrder", ctx, payment).
		Return("", expectedErr)

	// Act
	txID, err := svc.PayOrder(ctx, payment)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "", txID)
	assert.True(t, errors.Is(err, expectedErr))
	mockRepo.AssertExpectations(t)
}

func TestPayOrder_ValidationErrors(t *testing.T) {
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

func TestPayOrder_ContextHandling(t *testing.T) {
	testCases := []struct {
		name    string
		ctxFunc func() (context.Context, context.CancelFunc)
		expErr  error
		payment *model.Payment
	}{
		{
			name: "context canceled",
			ctxFunc: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx, cancel
			},
			expErr: context.Canceled,
			payment: &model.Payment{
				UserUuid:      "user-123",
				OrderUuid:     "order-123",
				PaymentMethod: model.CARD,
			},
		},
		{
			name: "context deadline exceeded",
			ctxFunc: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 0)
			},
			expErr: context.DeadlineExceeded,
			payment: &model.Payment{
				UserUuid:      "user-123",
				OrderUuid:     "order-123",
				PaymentMethod: model.CARD,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(mockRepo)
			svc := paymentService.NewService(mockRepo)

			ctx, cancel := tc.ctxFunc()
			defer cancel()

			mockRepo.On("PayOrder", ctx, tc.payment).Return("", tc.expErr)

			// Act
			txID, err := svc.PayOrder(ctx, tc.payment)

			// Assert
			assert.Error(t, err)
			assert.Equal(t, "", txID)
			assert.True(t, errors.Is(err, tc.expErr))
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetPayment_Success(t *testing.T) {
	// Arrange
	mockRepo := new(mockRepo)
	svc := paymentService.NewService(mockRepo)

	ctx := context.Background()
	expectedPayment := &model.Payment{
		UserUuid:      "user-123",
		OrderUuid:     "order-123",
		PaymentMethod: model.CARD,
	}
	txID := "tx-123"

	mockRepo.On("GetPayment", ctx, txID).Return(expectedPayment, nil)

	// Act
	payment, err := svc.GetPayment(ctx, txID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedPayment, payment)
	mockRepo.AssertExpectations(t)
}

func TestGetPayment_ValidationErrors(t *testing.T) {
	testCases := []struct {
		name   string
		txID   string
		expErr error
		setup  func(*mockRepo)
	}{
		{
			name:   "empty transaction ID",
			txID:   "",
			expErr: repoErrors.ErrTxIDRequired,
			setup: func(r *mockRepo) {
				r.On("GetPayment", mock.Anything, "").
					Return(nil, repoErrors.ErrTxIDRequired)
			},
		},
		{
			name:   "payment not found",
			txID:   "non-existent-tx",
			expErr: repoErrors.ErrPaymentNotFound,
			setup: func(r *mockRepo) {
				r.On("GetPayment", mock.Anything, "non-existent-tx").
					Return(nil, repoErrors.ErrPaymentNotFound)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(mockRepo)
			tc.setup(mockRepo)
			svc := paymentService.NewService(mockRepo)

			// Act
			payment, err := svc.GetPayment(context.Background(), tc.txID)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, payment)
			assert.True(t, errors.Is(err, tc.expErr))
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetPayment_ContextHandling(t *testing.T) {
	testCases := []struct {
		name    string
		ctxFunc func() (context.Context, context.CancelFunc)
		expErr  error
		txID    string
	}{
		{
			name: "context canceled",
			ctxFunc: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx, cancel
			},
			expErr: context.Canceled,
			txID:   "tx-123",
		},
		{
			name: "context deadline exceeded",
			ctxFunc: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 0)
			},
			expErr: context.DeadlineExceeded,
			txID:   "tx-123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(mockRepo)
			svc := paymentService.NewService(mockRepo)

			ctx, cancel := tc.ctxFunc()
			defer cancel()

			mockRepo.On("GetPayment", ctx, tc.txID).Return(nil, tc.expErr)

			// Act
			payment, err := svc.GetPayment(ctx, tc.txID)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, payment)
			assert.True(t, errors.Is(err, tc.expErr))
			mockRepo.AssertExpectations(t)
		})
	}
}
