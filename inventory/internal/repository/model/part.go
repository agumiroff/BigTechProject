package model

type Category int32

const (
	CategoryUnspecified Category = 0
	CategoryEngine      Category = 1
	CategoryFuel        Category = 2
	CategoryPorthole    Category = 3
	CategoryWing        Category = 4
)

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

// Value представляет тип oneof
type Value struct {
	StringValue *string
	Int64Value  *int64
}

type Part struct {
	Uuid          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    Dimensions
	Manufacturer  Manufacturer
	Tags          []string
	Metadata      map[string]*Value
	CreatedAt     int64
	UpdatedAt     int64
}

type PartsFilter struct {
	Uuids                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}
