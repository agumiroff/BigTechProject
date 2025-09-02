package repository

import (
	"context"
	"slices"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
	"github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/converter"
	rModel "github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/model"
)

func (s *repository) ListParts(ctx context.Context, filter *rModel.PartsFilter) ([]*model.Part, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*model.Part

	for _, part := range s.storage {
		if len(filter.Uuids) > 0 && !contains(filter.Uuids, part.Uuid) {
			continue
		}

		if len(filter.Names) > 0 && !contains(filter.Names, part.Name) {
			continue
		}

		if len(filter.Categories) > 0 && !containsCategory(filter.Categories, part.Category) {
			continue
		}

		if len(filter.ManufacturerCountries) > 0 && !contains(filter.ManufacturerCountries, part.Manufacturer.Country) {
			continue
		}

		if len(filter.Tags) > 0 && !hasAnyTag(part.Tags, filter.Tags) {
			continue
		}

		result = append(result, converter.RepoToModel(part))
	}
	return result, nil
}

func contains(list []string, value string) bool {
	return slices.Contains(list, value)
}

func containsCategory(list []rModel.Category, value rModel.Category) bool {
	return slices.Contains(list, value)
}

func hasAnyTag(partTags, filterTags []string) bool {
	tagSet := make(map[string]struct{}, len(partTags))
	for _, t := range partTags {
		tagSet[t] = struct{}{}
	}

	for _, tag := range filterTags {
		if _, found := tagSet[tag]; found {
			return true
		}
	}
	return false
}
