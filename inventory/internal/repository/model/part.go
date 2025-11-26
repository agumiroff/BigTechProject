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
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}

// Value — аналог protobuf oneof
type Value struct {
	StringValue *string  `bson:"string_value,omitempty"`
	Int64Value  *int64   `bson:"int64_value,omitempty"`
	BoolValue   *bool    `bson:"bool_value,omitempty"`
	DoubleValue *float64 `bson:"double_value,omitempty"`
}

type Part struct {
	UUID          string            `bson:"uuid"`
	Name          string            `bson:"name"`
	Description   string            `bson:"description"`
	Price         float64           `bson:"price"`
	StockQuantity int64             `bson:"stock_quantity"`
	Category      Category          `bson:"category"`
	Dimensions    Dimensions        `bson:"dimensions"`
	Manufacturer  Manufacturer      `bson:"manufacturer"`
	Tags          []string          `bson:"tags"`
	Metadata      map[string]*Value `bson:"metadata"`
	CreatedAt     int64             `bson:"created_at"`
	UpdatedAt     int64             `bson:"updated_at"`
}

type PartsFilter struct {
	UUIDs                 []string   `bson:"uuids,omitempty"`
	Names                 []string   `bson:"names,omitempty"`
	Categories            []Category `bson:"categories,omitempty"`
	ManufacturerCountries []string   `bson:"manufacturer_countries,omitempty"`
	Tags                  []string   `bson:"tags,omitempty"`
}
