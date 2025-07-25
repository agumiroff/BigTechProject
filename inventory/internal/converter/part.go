package converter

import (
	"github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
	invServiceV1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
)

func ModelToProto(m *model.Part) *invServiceV1.Part {
	if m == nil {
		return nil
	}

	return &invServiceV1.Part{
		Uuid:          m.Uuid,
		Name:          m.Name,
		Description:   m.Description,
		Price:         m.Price,
		StockQuantity: m.StockQuantity,
		Category:      invServiceV1.Category(m.Category),
		Dimensions: &invServiceV1.Dimensions{
			Length: m.Dimensions.Length,
			Width:  m.Dimensions.Width,
			Height: m.Dimensions.Height,
			Weight: m.Dimensions.Weight,
		},
		Manufacturer: &invServiceV1.Manufacturer{
			Name:    m.Manufacturer.Name,
			Country: m.Manufacturer.Country,
			Website: m.Manufacturer.Website,
		},
		Tags:      append([]string{}, m.Tags...),
		Metadata:  modelMetadataToProto(m.Metadata),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ProtoToModel(p *invServiceV1.Part) *model.Part {
	if p == nil {
		return nil
	}

	return &model.Part{
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
			Weight: p.Dimensions.Weight,
		},
		Manufacturer: model.Manufacturer{
			Name:    p.Manufacturer.Name,
			Country: p.Manufacturer.Country,
			Website: p.Manufacturer.Website,
		},
		Tags:      append([]string{}, p.Tags...),
		Metadata:  protoMetadataToModel(p.Metadata),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func modelMetadataToProto(meta map[string]*model.Value) map[string]*invServiceV1.Value {
	if meta == nil {
		return nil
	}

	result := make(map[string]*invServiceV1.Value, len(meta))
	for k, v := range meta {
		val := &invServiceV1.Value{}
		if v.StringValue != nil {
			val.Value = &invServiceV1.Value_StringValue{StringValue: *v.StringValue}
		} else if v.Int64Value != nil {
			val.Value = &invServiceV1.Value_Int64Value{Int64Value: *v.Int64Value}
		}
		result[k] = val
	}
	return result
}

func protoMetadataToModel(meta map[string]*invServiceV1.Value) map[string]*model.Value {
	if meta == nil {
		return nil
	}

	result := make(map[string]*model.Value, len(meta))
	for k, v := range meta {
		val := &model.Value{}
		switch x := v.Value.(type) {
		case *invServiceV1.Value_StringValue:
			val.StringValue = &x.StringValue
		case *invServiceV1.Value_Int64Value:
			val.Int64Value = &x.Int64Value
		}
		result[k] = val
	}
	return result
}
