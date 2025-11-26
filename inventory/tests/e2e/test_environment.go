package integration

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"go.mongodb.org/mongo-driver/bson"

	inventoryv1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
)

// InsertTestPart — вставляет тестовую запчасть в коллекцию Mongo и возвращает её UUID
func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partUUID := gofakeit.UUID()
	now := time.Now().Unix()

	partDoc := bson.M{
		"uuid":           partUUID,
		"name":           gofakeit.CarModel(),
		"description":    gofakeit.Sentence(6),
		"price":          gofakeit.Price(10, 5000),
		"stock_quantity": int64(gofakeit.IntRange(1, 100)),
		"category":       int32(inventoryv1.Category_CATEGORY_ENGINE),
		"dimensions": bson.M{
			"length": gofakeit.Float64Range(10, 500),
			"width":  gofakeit.Float64Range(10, 500),
			"height": gofakeit.Float64Range(10, 500),
			"weight": gofakeit.Float64Range(1, 50),
		},
		"manufacturer": bson.M{
			"name":    gofakeit.Company(),
			"country": gofakeit.Country(),
			"website": gofakeit.URL(),
		},
		"tags": []string{
			gofakeit.Word(),
			gofakeit.Word(),
		},
		"metadata": bson.M{
			"material": bson.M{
				"string_value": gofakeit.Word(),
			},
			"priority": bson.M{
				"int64_value": int64(gofakeit.IntRange(1, 100)),
			},
		},
		"created_at": now,
		"updated_at": now,
	}

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}

// InsertTestPartWithData — вставляет тестовую запчасть с заданными данными
func (env *TestEnvironment) InsertTestPartWithData(ctx context.Context, part *inventoryv1.Part) (string, error) {
	partDoc := bson.M{
		"uuid":           part.Uuid,
		"name":           part.Name,
		"description":    part.Description,
		"price":          part.Price,
		"stock_quantity": part.StockQuantity,
		"category":       int32(part.Category),
		"dimensions": bson.M{
			"length": part.Dimensions.Length,
			"width":  part.Dimensions.Width,
			"height": part.Dimensions.Height,
			"weight": part.Dimensions.Weight,
		},
		"manufacturer": bson.M{
			"name":    part.Manufacturer.Name,
			"country": part.Manufacturer.Country,
			"website": part.Manufacturer.Website,
		},
		"tags":       part.Tags,
		"metadata":   convertMetadataToBson(part.Metadata),
		"created_at": part.CreatedAt,
		"updated_at": part.UpdatedAt,
	}

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return part.Uuid, nil
}

// convertMetadataToBson конвертирует protobuf metadata в BSON формат
func convertMetadataToBson(metadata map[string]*inventoryv1.Value) bson.M {
	result := bson.M{}
	for key, value := range metadata {
		switch v := value.Value.(type) {
		case *inventoryv1.Value_StringValue:
			result[key] = bson.M{"string_value": v.StringValue}
		case *inventoryv1.Value_Int64Value:
			result[key] = bson.M{"int64_value": v.Int64Value}
		case *inventoryv1.Value_BoolValue:
			result[key] = bson.M{"bool_value": v.BoolValue}
		case *inventoryv1.Value_DoubleValue:
			result[key] = bson.M{"double_value": v.DoubleValue}
		}
	}
	return result
}

// GetTestPartData — возвращает тестовую запчасть для создания (без UUID и timestamps)
func (env *TestEnvironment) GetTestPartData() *inventoryv1.Part {
	return &inventoryv1.Part{
		Name:          gofakeit.CarModel(),
		Description:   gofakeit.Sentence(6),
		Price:         gofakeit.Price(10, 5000),
		StockQuantity: int64(gofakeit.IntRange(1, 100)),
		Category:      inventoryv1.Category_CATEGORY_ENGINE,
		Dimensions: &inventoryv1.Dimensions{
			Length: gofakeit.Float64Range(10, 500),
			Width:  gofakeit.Float64Range(10, 500),
			Height: gofakeit.Float64Range(10, 500),
			Weight: gofakeit.Float64Range(1, 50),
		},
		Manufacturer: &inventoryv1.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags: []string{
			gofakeit.Word(),
			gofakeit.Word(),
		},
		Metadata: map[string]*inventoryv1.Value{
			"material": {
				Value: &inventoryv1.Value_StringValue{
					StringValue: gofakeit.Word(),
				},
			},
			"priority": {
				Value: &inventoryv1.Value_Int64Value{
					Int64Value: int64(gofakeit.IntRange(1, 100)),
				},
			},
		},
	}
}

// CreateTestPart — создает полную тестовую запчасть (с UUID и timestamps)
func (env *TestEnvironment) CreateTestPart() *inventoryv1.Part {
	now := time.Now().Unix()
	part := env.GetTestPartData()
	part.Uuid = gofakeit.UUID()
	part.CreatedAt = now
	part.UpdatedAt = now

	return part
}

// ClearInventoryCollection — удаляет все записи из коллекции parts
func (env *TestEnvironment) ClearInventoryCollection(ctx context.Context) error {
	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
