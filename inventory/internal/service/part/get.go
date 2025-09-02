package part

import (
	"context"
	"log"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
)

func (s *service) GetPart(ctx context.Context, uuid string) (res *model.Part, err error) {
	m, err := s.Repo.GetPart(ctx, uuid)
	if err != nil {
		log.Printf("failed to get part: %v", err)
		return &model.Part{}, err
	}

	return m, nil
}
