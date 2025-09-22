package converter

import (
	"github.com/agumiroff/BigTechProject/order/v1/external/repository/model"
	inventoryv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/payment/v1"
)

func ModelPaymentToProto(m *model.Payment) *paymentv1.Payment {
	return &paymentv1.Payment{
		OrderUuid:     m.OrderUUID,
		UserUuid:      m.UserUUID,
		PaymentMethod: mapPaymentMethod(m.PaymentMethod),
	}
}

func ProtoPartToModel(protoParts []*inventoryv1.Part) []*model.Part {
	var result []*model.Part

	for _, p := range protoParts {
		part := &model.Part{
			Uuid:          p.Uuid,
			Name:          p.Name,
			Description:   p.Description,
			Price:         p.Price,
			StockQuantity: p.StockQuantity,
			Category:      int32(p.Category),
			Dimensions: model.Dimensions{
				Length: p.Dimensions.Length,
				Width:  p.Dimensions.Width,
				Height: p.Dimensions.Height,
			},
			Manufacturer: model.Manufacturer{
				Name:    p.Manufacturer.Name,
				Country: p.Manufacturer.Country,
			},
			Tags:      p.Tags,
			Metadata:  convertMetadataToModel(p.Metadata),
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
		result = append(result, part)
	}

	return result
}

func convertMetadataToModel(meta map[string]*inventoryv1.Value) map[string]*model.Value {
	result := make(map[string]*model.Value)
	for k, v := range meta {
		m := &model.Value{}

		switch val := v.Value.(type) {
		case *inventoryv1.Value_StringValue:
			m.StringValue = &val.StringValue
		case *inventoryv1.Value_Int64Value:
			m.Int64Value = &val.Int64Value
		}

		result[k] = m
	}
	return result
}

func ModelFilterToProto(m *model.PartsFilter) *inventoryv1.PartsFilter {
	return &inventoryv1.PartsFilter{
		Uuids:                 m.Uuids,
		Names:                 m.Names,
		ManufacturerCountries: m.ManufacturerCountries,
		Categories:            mapCategories(m.Categories),
	}
}

func mapCategories(c []model.Category) []inventoryv1.Category {
	var result []inventoryv1.Category
	for _, cat := range c {
		switch cat {
		case model.CategoryEngine:
			result = append(result, inventoryv1.Category_CATEGORY_ENGINE)
		case model.CategoryFuel:
			result = append(result, inventoryv1.Category_CATEGORY_FUEL)
		case model.CategoryWing:
			result = append(result, inventoryv1.Category_CATEGORY_WING)
		case model.CategoryPorthole:
			result = append(result, inventoryv1.Category_CATEGORY_PORTHOLE)
		default:
			result = append(result, inventoryv1.Category_CATEGORY_UNSPECIFIED)
		}
	}

	return result
}

func ProtoListToModel(p *inventoryv1.ListPartsRequest) *model.ListPartsRequest {
	return &model.ListPartsRequest{
		PartsFilter: model.PartsFilter{
			Uuids:                 p.GetFilter().GetUuids(),
			Names:                 p.GetFilter().GetNames(),
			ManufacturerCountries: p.GetFilter().GetManufacturerCountries(),
			Categories:            mapProtoCategories(p.GetFilter().GetCategories()),
		},
	}
}

func mapProtoCategories(c []inventoryv1.Category) []model.Category {
	var result []model.Category
	for _, cat := range c {
		switch cat {
		case inventoryv1.Category_CATEGORY_ENGINE:
			result = append(result, model.CategoryEngine)
		case inventoryv1.Category_CATEGORY_FUEL:
			result = append(result, model.CategoryFuel)
		case inventoryv1.Category_CATEGORY_WING:
			result = append(result, model.CategoryWing)
		case inventoryv1.Category_CATEGORY_PORTHOLE:
			result = append(result, model.CategoryPorthole)
		default:
			result = append(result, model.CategoryUnspecified)
		}
	}
	return result
}

func mapPaymentMethod(method model.PaymentMethod) paymentv1.PaymentMethod {
	switch method {
	case model.PaymentMethodCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCREDITCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodINVESTORMONEY:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
