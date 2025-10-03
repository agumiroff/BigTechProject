package part

import (
	"context"
	"errors"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
)

func (s *service) GetPart(ctx context.Context, uuid string) (res *model.Part, err error) {
	if uuid == "" {
		return &model.Part{}, errors.New("uuid is empty")
	}

	m, err := s.Repo.GetPart(ctx, uuid)
	if err != nil {
		return &model.Part{}, err
	}

	return m, nil
}
