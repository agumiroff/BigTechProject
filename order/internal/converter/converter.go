package converter

import (
	"github.com/agumiroff/BigTechProject/order/v1/internal/model"
	rModel "github.com/agumiroff/BigTechProject/order/v1/internal/repository/model"
	OrderV1 "github.com/agumiroff/BigTechProject/shared/pkg/openapi/v1"
	paymentv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

func ToModelCreateOrderRequest(req *OrderV1.CreateOrderRequest) *model.CreateOrderRequest {
	return &model.CreateOrderRequest{
		UserUUID:  req.GetUserUUID(),
		PartUUIDs: req.GetPartUuids(),
	}
}

func ToRepoOrder(m *model.Order) *rModel.Order {
	return &rModel.Order{
		OrderUUID:       m.OrderUUID,
		UserUUID:        m.UserUUID,
		PartUUIDs:       m.PartUUIDs,
		TotalPrice:      m.TotalPrice,
		TransactionUUID: m.TransactionUUID,
		PaymentMethod:   rModel.PaymentMethod(m.PaymentMethod),
		Status:          rModel.OrderStatus(m.Status),
	}
}

func ToProtoOrder(m *model.Order) *OrderV1.Order {
	return &OrderV1.Order{
		OrderUUID:       m.OrderUUID,
		UserUUID:        m.UserUUID,
		PartUuids:       m.PartUUIDs,
		TotalPrice:      m.TotalPrice,
		TransactionUUID: OrderV1.NewOptNilString(m.TransactionUUID),
		PaymentMethod:   OrderV1.NewOptPaymentMethod(OrderV1.PaymentMethod(m.PaymentMethod)),
		Status:          OrderV1.OrderStatus(m.Status),
	}
}

func ToModelOrder(m *rModel.Order) *model.Order {
	return &model.Order{
		OrderUUID:       m.OrderUUID,
		UserUUID:        m.UserUUID,
		PartUUIDs:       m.PartUUIDs,
		TotalPrice:      m.TotalPrice,
		TransactionUUID: m.TransactionUUID,
		PaymentMethod:   model.PaymentMethod(m.PaymentMethod),
		Status:          model.OrderStatus(m.Status),
	}
}

func ToProtoPaymentMethod(m *model.PaymentMethod) paymentv1.PaymentMethod {
	switch *m {
	case model.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodINVESTORMONEY:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

func ToModelPaymentMethod(method OrderV1.PaymentMethod) model.PaymentMethod {
	switch method {
	case OrderV1.PaymentMethodCARD:
		return model.PaymentMethodCARD
	case OrderV1.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case OrderV1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCREDITCARD
	case OrderV1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodINVESTORMONEY
	default:
		return model.PaymentMethodUNKNOWN
	}
}

func ToModelOrderFromProto(m *OrderV1.Order) *model.Order {
	return &model.Order{
		OrderUUID:       m.OrderUUID,
		UserUUID:        m.UserUUID,
		PartUUIDs:       m.PartUuids,
		TotalPrice:      m.TotalPrice,
		TransactionUUID: m.TransactionUUID.Value,
		PaymentMethod:   ToModelPaymentMethod(m.PaymentMethod.Value),
		Status:          model.OrderStatus(m.Status),
	}
}
