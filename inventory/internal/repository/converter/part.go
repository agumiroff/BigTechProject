package converter

import (
	"log"

	model "github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
	rModel "github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/model"
)

func ModelToRepo(m *model.Part) *rModel.Part {
	if m == nil {
		log.Printf("Business model is empty: %v", m)
		return nil
	}

	return &rModel.Part{
		Uuid:          m.Uuid,
		Name:          m.Name,
		Description:   m.Description,
		Price:         m.Price,
		StockQuantity: m.StockQuantity,
		Category:      rModel.Category(m.Category),
		Dimensions: rModel.Dimensions{
			Length: m.Dimensions.Length,
			Width:  m.Dimensions.Width,
			Height: m.Dimensions.Height,
			Weight: m.Dimensions.Weight,
		},
		Manufacturer: rModel.Manufacturer{
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

func RepoToModel(m *rModel.Part) *model.Part {
	if m == nil {
		log.Printf("Repo model is empty :%v", m)
		return nil
	}

	return &model.Part{
		Uuid:          m.Uuid,
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
func FilterToRepo(f *model.PartsFilter) *rModel.PartsFilter {
	return &rModel.PartsFilter{
		Uuids:                 f.Uuids,
		Names:                 f.Names,
		Categories:            CategoryToRepo(f.Categories),
		ManufacturerCountries: f.ManufacturerCountries,
	}
}

// CategoryToRepo Category convert
func CategoryToRepo(c []model.Category) []rModel.Category {
	result := make([]rModel.Category, 0, len(c))
	for _, c := range c {
		result = append(result, categoryToRepo(c))
	}
	return result
}

func categoryToRepo(c model.Category) rModel.Category {
	switch c {
	case model.CategoryEngine:
		return rModel.CategoryEngine
	case model.CategoryFuel:
		return rModel.CategoryFuel
	case model.CategoryPorthole:
		return rModel.CategoryPorthole
	case model.CategoryWing:
		return rModel.CategoryWing
	default:
		return rModel.CategoryUnspecified
	}
}

// Metadata converting
func cloneMetadataToRepo(in map[string]*model.Value) map[string]*rModel.Value {
	if in == nil {
		return nil
	}
	out := make(map[string]*rModel.Value, len(in))
	for k, v := range in {
		val := &rModel.Value{}
		if v.StringValue != nil {
			s := *v.StringValue
			val.StringValue = &s
		}
		if v.Int64Value != nil {
			i := *v.Int64Value
			val.Int64Value = &i
		}

		out[k] = val
	}
	return out
}

func cloneMetadata(in map[string]*rModel.Value) map[string]*model.Value {
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
		out[k] = val
	}
	return out
}
