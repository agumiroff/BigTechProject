package converter

import (
	"github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
	inventoryv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
)

func ToProtoModels(parts []*model.Part) []*inventoryv1.Part {
	var protos []*inventoryv1.Part
	for _, part := range parts {
		protos = append(protos, ToProtoPart(part))
	}
	return protos
}

func ToFilterProto(f *model.PartsFilter) *inventoryv1.PartsFilter {
	return &inventoryv1.PartsFilter{
		Uuids:                 f.Uuids,
		Names:                 f.Names,
		Categories:            ConvertModelCategoriesToProto(f.Categories),
		ManufacturerCountries: f.ManufacturerCountries,
		Tags:                  f.Tags,
	}
}

func ToFilterModel(f *inventoryv1.PartsFilter) *model.PartsFilter {
	return &model.PartsFilter{
		Uuids:                 f.Uuids,
		Names:                 f.Names,
		Categories:            ConvertProtoCategoriesToModel(f.Categories),
		ManufacturerCountries: f.ManufacturerCountries,
		Tags:                  f.Tags,
	}
}

func ConvertModelCategoriesToProto(categories []model.Category) []inventoryv1.Category {
	result := make([]inventoryv1.Category, 0, len(categories))
	for _, c := range categories {
		result = append(result, ConvertModelToProtoCategory(c))
	}
	return result
}

func ConvertModelToProtoCategory(c model.Category) inventoryv1.Category {
	switch c {
	case model.CategoryUnspecified:
		return inventoryv1.Category_CATEGORY_UNSPECIFIED
	case model.CategoryEngine:
		return inventoryv1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return inventoryv1.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return inventoryv1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return inventoryv1.Category_CATEGORY_WING
	default:
		return inventoryv1.Category_CATEGORY_UNSPECIFIED
	}
}

func ConvertProtoCategoriesToModel(categories []inventoryv1.Category) []model.Category {
	result := make([]model.Category, 0, len(categories))
	for _, c := range categories {
		result = append(result, ConvertProtoToModelCategory(c))
	}
	return result
}

func ConvertProtoToModelCategory(c inventoryv1.Category) model.Category {
	switch c {
	case inventoryv1.Category_CATEGORY_ENGINE:
		return model.CategoryEngine
	case inventoryv1.Category_CATEGORY_FUEL:
		return model.CategoryFuel
	case inventoryv1.Category_CATEGORY_PORTHOLE:
		return model.CategoryPorthole
	case inventoryv1.Category_CATEGORY_WING:
		return model.CategoryWing
	case inventoryv1.Category_CATEGORY_UNSPECIFIED:
		fallthrough
	default:
		return model.CategoryUnspecified
	}
}

func ToProtoPart(m *model.Part) *inventoryv1.Part {
	if m == nil {
		return nil
	}

	return &inventoryv1.Part{
		Uuid:          m.Uuid,
		Name:          m.Name,
		Description:   m.Description,
		Price:         m.Price,
		StockQuantity: m.StockQuantity,
		Category:      inventoryv1.Category(m.Category),
		Dimensions: &inventoryv1.Dimensions{
			Length: m.Dimensions.Length,
			Width:  m.Dimensions.Width,
			Height: m.Dimensions.Height,
			Weight: m.Dimensions.Weight,
		},
		Manufacturer: &inventoryv1.Manufacturer{
			Name:    m.Manufacturer.Name,
			Country: m.Manufacturer.Country,
			Website: m.Manufacturer.Website,
		},
		Tags:      append([]string{}, m.Tags...),
		Metadata:  ToProtoMetadata(m.Metadata),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ToModelPart(p *inventoryv1.Part) *model.Part {
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
		Metadata:  ToModelMetadata(p.Metadata),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func ToProtoMetadata(meta map[string]*model.Value) map[string]*inventoryv1.Value {
	if meta == nil {
		return nil
	}

	result := make(map[string]*inventoryv1.Value, len(meta))
	for k, v := range meta {
		val := &inventoryv1.Value{}
		if v.StringValue != nil {
			val.Value = &inventoryv1.Value_StringValue{StringValue: *v.StringValue}
		} else if v.Int64Value != nil {
			val.Value = &inventoryv1.Value_Int64Value{Int64Value: *v.Int64Value}
		}
		result[k] = val
	}
	return result
}

func ToModelMetadata(meta map[string]*inventoryv1.Value) map[string]*model.Value {
	if meta == nil {
		return nil
	}

	result := make(map[string]*model.Value, len(meta))
	for k, v := range meta {
		val := &model.Value{}
		switch x := v.Value.(type) {
		case *inventoryv1.Value_StringValue:
			val.StringValue = &x.StringValue
		case *inventoryv1.Value_Int64Value:
			val.Int64Value = &x.Int64Value
		}
		result[k] = val
	}
	return result
}
