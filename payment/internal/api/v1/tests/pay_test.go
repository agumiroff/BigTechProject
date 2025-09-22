package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	v1 "github.com/agumiroff/BigTechProject/payment/v1/internal/api/v1"
	"github.com/agumiroff/BigTechProject/payment/v1/internal/model"
	repoErrors "github.com/agumiroff/BigTechProject/payment/v1/internal/repository/payment"
	paymentv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) PayOrder(ctx context.Context, p *model.Payment) (string, error) {
	args := m.Called(ctx, p)
	return args.String(0), args.Error(1)
}

func (m *mockService) GetPayment(ctx context.Context, uuid string) (*model.Payment, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Payment), args.Error(1)
}

func TestPayOrder_Success(t *testing.T) {
	// Arrange
	mockSvc := new(mockService)
	api := v1.NewAPI(mockSvc)

	ctx := context.Background()
	req := &paymentv1.PayOrderRequest{
		Payment: &paymentv1.Payment{
			UserUuid:      "user-123",
			OrderUuid:     "order-123",
			PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
		},
	}

	expectedTxID := "tx-123"
	mockSvc.On("PayOrder", ctx, &model.Payment{
		UserUuid:      "user-123",
		OrderUuid:     "order-123",
		PaymentMethod: model.CARD,
	}).Return(expectedTxID, nil)

	// Act
	resp, err := api.PayOrder(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedTxID, resp.TransactionUuid)
	mockSvc.AssertExpectations(t)
}

func TestPayOrder_ValidationErrors(t *testing.T) {
	testCases := []struct {
		name     string
		req      *paymentv1.PayOrderRequest
		mockFunc func(*mockService)
		expErr   error
	}{
		{
			name: "nil payment",
			req:  &paymentv1.PayOrderRequest{Payment: nil},
			mockFunc: func(s *mockService) {
				s.On("PayOrder", mock.Anything, (*model.Payment)(nil)).
					Return("", repoErrors.ErrPaymentRequired)
			},
			expErr: repoErrors.ErrPaymentRequired,
		},
		{
			name: "empty user uuid",
			req: &paymentv1.PayOrderRequest{
				Payment: &paymentv1.Payment{
					OrderUuid:     "order-123",
					PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
				},
			},
			mockFunc: func(s *mockService) {
				s.On("PayOrder", mock.Anything, mock.MatchedBy(func(p *model.Payment) bool {
					return p.OrderUuid == "order-123" && p.UserUuid == "" && p.PaymentMethod == model.CARD
				})).Return("", repoErrors.ErrUserUUIDRequired)
			},
			expErr: repoErrors.ErrUserUUIDRequired,
		},
		{
			name: "empty order uuid",
			req: &paymentv1.PayOrderRequest{
				Payment: &paymentv1.Payment{
					UserUuid:      "user-123",
					PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
				},
			},
			mockFunc: func(s *mockService) {
				s.On("PayOrder", mock.Anything, mock.MatchedBy(func(p *model.Payment) bool {
					return p.UserUuid == "user-123" && p.OrderUuid == "" && p.PaymentMethod == model.CARD
				})).Return("", repoErrors.ErrOrderUUIDRequired)
			},
			expErr: repoErrors.ErrOrderUUIDRequired,
		},
		{
			name: "invalid payment method",
			req: &paymentv1.PayOrderRequest{
				Payment: &paymentv1.Payment{
					UserUuid:      "user-123",
					OrderUuid:     "order-123",
					PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
				},
			},
			mockFunc: func(s *mockService) {
				s.On("PayOrder", mock.Anything, mock.MatchedBy(func(p *model.Payment) bool {
					return p.UserUuid == "user-123" && p.OrderUuid == "order-123" && p.PaymentMethod == model.CategoryUnspecified
				})).Return("", repoErrors.ErrPaymentMethodInvalid)
			},
			expErr: repoErrors.ErrPaymentMethodInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockSvc := new(mockService)
			tc.mockFunc(mockSvc)
			api := v1.NewAPI(mockSvc)

			// Act
			resp, err := api.PayOrder(context.Background(), tc.req)

			// Assert
			assert.Error(t, err)
			assert.NotNil(t, resp)
			assert.Empty(t, resp.TransactionUuid)
			assert.True(t, errors.Is(err, tc.expErr))
			mockSvc.AssertExpectations(t)
		})
	}
}
