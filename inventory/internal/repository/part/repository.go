package repository

import (
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/samber/lo"

	rep "github.com/agumiroff/BigTechProject/inventory/v1/internal/repository"
	model "github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/model"
)

type repository struct {
	mu      sync.RWMutex
	storage map[string]*model.Part
}

var _ rep.InvRepository = (*repository)(nil)

func NewRepository() (res *repository) {
	service := &repository{
		storage: make(map[string]*model.Part),
	}

	fillStorage(service)
	return service
}

func fillStorage(s *repository) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tags := []string{
		gofakeit.Word(),
		gofakeit.BuzzWord(),
	}

	metadata := map[string]*model.Value{
		"color": {
			StringValue: lo.ToPtr(gofakeit.Color()),
		},
		"year": {
			Int64Value: lo.ToPtr(int64(gofakeit.Year())),
		},
	}

	now := time.Now().Unix()

	iterations := 10
	for range iterations {
		part := &model.Part{
			Uuid:          gofakeit.UUID(),
			Name:          gofakeit.Name(),
			Description:   gofakeit.HipsterSentence(5),
			Price:         gofakeit.Price(10, 1000),
			StockQuantity: int64(gofakeit.Number(1, 500)),
			Category:      model.CategoryEngine, // #nosec G115
			Dimensions: model.Dimensions{
				Length: gofakeit.Float64Range(10, 200),
				Width:  gofakeit.Float64Range(10, 200),
				Height: gofakeit.Float64Range(10, 200),
				Weight: gofakeit.Float64Range(1, 100),
			},
			Manufacturer: model.Manufacturer{
				Name:    gofakeit.Company(),
				Country: gofakeit.Country(),
				Website: gofakeit.URL(),
			},
			Tags:      tags,
			Metadata:  metadata,
			CreatedAt: now - 86400,
			UpdatedAt: now,
		}

		s.storage[part.Uuid] = part
	}
}
