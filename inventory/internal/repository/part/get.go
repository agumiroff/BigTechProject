package repository

import (
	"context"
	"log"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
	"github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/converter"
)

func (s *repository) GetPart(ctx context.Context, uuid string) (res *model.Part, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if uuid == "" {
		log.Printf("error %d\n", err)
		return nil, err
	}

	part := s.storage[uuid]

	converted := converter.RepoToModel(part)

	return converted, nil
}
