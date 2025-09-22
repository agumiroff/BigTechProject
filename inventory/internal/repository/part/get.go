package repository

import (
	"context"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
	"github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/converter"
)

func (s *repository) GetPart(ctx context.Context, uuid string) (res *model.Part, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part := s.storage[uuid]

	convertedPart := converter.RepoToModel(part)

	return convertedPart, nil
}
