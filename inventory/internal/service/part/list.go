package part

import (
	"context"
	"github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
	rModel "github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/model"
	"log"
)

func (s *service) LitParts(ctx context.Context, f *rModel.PartsFilter) (res []*model.Part, err error) {
	part, err := s.Repo.ListParts(ctx, f)
	if err != nil {
		log.Printf("Part not found: %d", err)
		return []*model.Part{}, err
	}

	return part, nil
}
