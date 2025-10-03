package part

import (
	"context"
	"log"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
	rConverter "github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/converter"
)

func (s *service) ListParts(ctx context.Context, f *model.PartsFilter) (res []*model.Part, err error) {
	part, err := s.Repo.ListParts(ctx, rConverter.FilterToRepo(f))
	if err != nil {
		log.Printf("part not found: %d", err)
		return []*model.Part{}, err
	}

	return part, nil
}
