package part

import (
	"context"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
)

func (s *service) CreatePart(ctx context.Context, p *model.Part) (res string, err error) {
	uuid, err := s.Repo.CreatePart(ctx, p)
	if err != nil {
		return "", err
	}

	return uuid, nil
}
