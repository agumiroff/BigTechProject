package model

type ListPartsRequest struct {
	PartsFilter PartsFilter
}

type ListPartsResponse struct {
	Parts []*Part
}

type Category int32

const (
	CategoryUnspecified Category = 0
	CategoryEngine      Category = 1
	CategoryFuel        Category = 2
	CategoryPorthole    Category = 3
	CategoryWing        Category = 4
)

func (c Category) String() string {
	switch c {
	case CategoryUnspecified:
		return "CATEGORY_UNSPECIFIED"
	case CategoryEngine:
		return "CATEGORY_ENGINE"
	case CategoryFuel:
		return "CATEGORY_FUEL"
	case CategoryPorthole:
		return "CATEGORY_PORTHOLE"
	case CategoryWing:
		return "CATEGORY_WING"
	default:
		return "UNKNOWN_CATEGORY"
	}
}

type Part struct {
	Uuid          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      int32
	Dimensions    Dimensions
	Manufacturer  Manufacturer
	Tags          []string
	Metadata      map[string]*Value
	CreatedAt     int64
	UpdatedAt     int64
}

type PartsFilter struct {
	Uuids                 []string   `json:"uuids"`
	Names                 []string   `json:"names"`
	Categories            []Category `json:"categories"`
	ManufacturerCountries []string   `json:"manufacturer_countries"`
	Tags                  []string   `json:"tags"`
}

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type Value struct {
	StringValue *string
	Int64Value  *int64
}
