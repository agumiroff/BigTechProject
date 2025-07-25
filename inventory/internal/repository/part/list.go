package part

import (
	"context"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
	"github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/converter"
	rModel "github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/model"
)

func (s *repository) ListParts(ctx context.Context, filter *rModel.PartsFilter) ([]*model.Part, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*model.Part

	for _, part := range s.storage {
		// Фильтрация по UUID
		if len(filter.Uuids) > 0 && !contains(filter.Uuids, part.Uuid) {
			continue
		}

		// Фильтрация по Name
		if len(filter.Names) > 0 && !contains(filter.Names, part.Name) {
			continue
		}

		// Фильтрация по Category
		if len(filter.Categories) > 0 && !containsCategory(filter.Categories, part.Category) {
			continue
		}

		// Фильтрация по Manufacturer.Country
		if len(filter.ManufacturerCountries) > 0 && !contains(filter.ManufacturerCountries, part.Manufacturer.Country) {
			continue
		}

		// Фильтрация по Tags (логическое ИЛИ: хотя бы один тег должен совпасть)
		if len(filter.Tags) > 0 && !hasAnyTag(part.Tags, filter.Tags) {
			continue
		}

		result = append(result, converter.RepoToModel(part))
	}
	return result, nil
}

func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func containsCategory(list []rModel.Category, value rModel.Category) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
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
