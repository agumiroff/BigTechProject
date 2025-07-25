package converter

import (
	"log"

	dModel "github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
	rModel "github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/model"
)

func ModelToRepo(m *dModel.Part) *rModel.Part {
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

func RepoToModel(m *rModel.Part) *dModel.Part {
	if m == nil {
		log.Printf("Repo model is empty :%v", m)
		return nil
	}

	return &dModel.Part{
		Uuid:          m.Uuid,
		Name:          m.Name,
		Description:   m.Description,
		Price:         m.Price,
		StockQuantity: m.StockQuantity,
		Category:      int32(m.Category),
		Dimensions: dModel.Dimensions{
			Length: m.Dimensions.Length,
			Width:  m.Dimensions.Width,
			Height: m.Dimensions.Height,
			Weight: m.Dimensions.Weight,
		},
		Manufacturer: dModel.Manufacturer{
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

func cloneMetadataToRepo(in map[string]*dModel.Value) map[string]*rModel.Value {
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

func cloneMetadata(in map[string]*rModel.Value) map[string]*dModel.Value {
	if in == nil {
		return nil
	}
	out := make(map[string]*dModel.Value, len(in))
	for k, v := range in {
		val := &dModel.Value{}
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
