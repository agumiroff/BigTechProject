package part

import (
	"context"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func (s *service) GetPart(ctx context.Context, uuid string) (res *model.Part, err error) {
	m, err := s.Repo.GetPart(ctx, uuid)

	if uuid == " " {
		return &model.Part{}, apperrors.ErrInvalidRequest
	}

	if err != nil {
		return &model.Part{}, apperrors.ErrInvalidRequest
	}

	return m, nil
}
