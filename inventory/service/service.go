package service

import (
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v6"

	invServiceV1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
)

type InvService struct {
	invServiceV1.UnimplementedInventoryServiceServer

	mu      sync.RWMutex
	storage map[string]*invServiceV1.Part
}

func NewService() (res *InvService) {
	service := &InvService{
		storage: make(map[string]*invServiceV1.Part),
	}

	fillStorage(service)
	return service
}

func fillStorage(s *InvService) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tags := []string{
		gofakeit.Word(),
		gofakeit.BuzzWord(),
	}

	metadata := map[string]*invServiceV1.Value{
		"color": {Value: &invServiceV1.Value_StringValue{StringValue: gofakeit.Color()}},
		"year":  {Value: &invServiceV1.Value_Int64Value{Int64Value: int64(gofakeit.Year())}},
	}

	now := time.Now().Unix()

	for i := 0; i < 10; i++ {
		part := &invServiceV1.Part{
			Uuid:          gofakeit.UUID(),
			Name:          gofakeit.Name(),
			Description:   gofakeit.HipsterSentence(5),
			Price:         gofakeit.Price(10, 1000),
			StockQuantity: int64(gofakeit.Number(1, 500)),
			Category:      invServiceV1.Category(1),
			Dimensions: &invServiceV1.Dimensions{
				Length: gofakeit.Float64Range(10, 200),
				Width:  gofakeit.Float64Range(10, 200),
				Height: gofakeit.Float64Range(10, 200),
				Weight: gofakeit.Float64Range(1, 100),
			},
			Manufacturer: &invServiceV1.Manufacturer{
				Name:    gofakeit.Company(),
				Country: gofakeit.Country(),
				Website: gofakeit.URL(),
			},
			Tags:      tags,
			Metadata:  metadata,
			CreatedAt: now - 86400, // сутки назад
			UpdatedAt: now,
		}

		s.storage[part.GetUuid()] = part
	}
}
