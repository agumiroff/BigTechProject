package converter

import (
	model "github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
	repomodel "github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/model"
)

func ModelToRepo(m *model.Part) *repomodel.Part {
	if m == nil {
		return nil
	}

	return &repomodel.Part{
		UUID:          m.Uuid,
		Name:          m.Name,
		Description:   m.Description,
		Price:         m.Price,
		StockQuantity: m.StockQuantity,
		Category:      repomodel.Category(m.Category),
		Dimensions: repomodel.Dimensions{
			Length: m.Dimensions.Length,
			Width:  m.Dimensions.Width,
			Height: m.Dimensions.Height,
			Weight: m.Dimensions.Weight,
		},
		Manufacturer: repomodel.Manufacturer{
			Name:    m.Manufacturer.Name,
			Country: m.Manufacturer.Country,
			Website: m.Manufacturer.Website,
		},
		Tags:      append([]string{}, m.Tags...),
		Metadata:  cloneMetadataToRepo(m.Metadata),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func RepoToModel(m *repomodel.Part) *model.Part {
	if m == nil {
		return nil
	}

	return &model.Part{
		Uuid:          m.UUID,
		Name:          m.Name,
		Description:   m.Description,
		Price:         m.Price,
		StockQuantity: m.StockQuantity,
		Category:      int32(m.Category),
		Dimensions: model.Dimensions{
			Length: m.Dimensions.Length,
			Width:  m.Dimensions.Width,
			Height: m.Dimensions.Height,
			Weight: m.Dimensions.Weight,
		},
		Manufacturer: model.Manufacturer{
			Name:    m.Manufacturer.Name,
			Country: m.Manufacturer.Country,
			Website: m.Manufacturer.Website,
		},
		Tags:      append([]string{}, m.Tags...), // копируем срез
		Metadata:  cloneMetadata(m.Metadata),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// FilterToRepo Filter convert
func FilterToRepo(f *model.PartsFilter) *repomodel.PartsFilter {
	return &repomodel.PartsFilter{
		UUIDs:                 f.Uuids,
		Names:                 f.Names,
		Categories:            CategoryToRepo(f.Categories),
		ManufacturerCountries: f.ManufacturerCountries,
		Tags:                  f.Tags,
	}
}

// CategoryToRepo Category convert
func CategoryToRepo(c []model.Category) []repomodel.Category {
	result := make([]repomodel.Category, 0, len(c))
	for _, c := range c {
		result = append(result, categoryToRepo(c))
	}
	return result
}

func categoryToRepo(c model.Category) repomodel.Category {
	switch c {
	case model.CategoryEngine:
		return repomodel.CategoryEngine
	case model.CategoryFuel:
		return repomodel.CategoryFuel
	case model.CategoryPorthole:
		return repomodel.CategoryPorthole
	case model.CategoryWing:
		return repomodel.CategoryWing
	default:
		return repomodel.CategoryUnspecified
	}
}

// Metadata converting
func cloneMetadataToRepo(in map[string]*model.Value) map[string]*repomodel.Value {
	if in == nil {
		return nil
	}
	out := make(map[string]*repomodel.Value, len(in))
	for k, v := range in {
		val := &repomodel.Value{}
		if v.StringValue != nil {
			s := *v.StringValue
			val.StringValue = &s
		}
		if v.Int64Value != nil {
			i := *v.Int64Value
			val.Int64Value = &i
		}
		if v.BoolValue != nil {
			b := *v.BoolValue
			val.BoolValue = &b
		}
		if v.DoubleValue != nil {
			d := *v.DoubleValue
			val.DoubleValue = &d
		}
		out[k] = val
	}
	return out
}

func cloneMetadata(in map[string]*repomodel.Value) map[string]*model.Value {
	if in == nil {
		return nil
	}
	out := make(map[string]*model.Value, len(in))
	for k, v := range in {
		val := &model.Value{}
		if v.StringValue != nil {
			s := *v.StringValue
			val.StringValue = &s
		}
		if v.Int64Value != nil {
			i := *v.Int64Value
			val.Int64Value = &i
		}
		if v.BoolValue != nil {
			b := *v.BoolValue
			val.BoolValue = &b
		}
		if v.DoubleValue != nil {
			d := *v.DoubleValue
			val.DoubleValue = &d
		}
		out[k] = val
	}
	return out
}
