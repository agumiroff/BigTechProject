package service

import (
	"context"
	"log"

	invServiceV1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
)

func (s *InvService) GetPart(ctx context.Context, req *invServiceV1.GetPartRequest) (res *invServiceV1.GetPartResponse, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id := req.GetUuid()
	if id == "" {
		log.Printf("error %d\n", err)
		return nil, err
	}

	part := &invServiceV1.GetPartResponse{
		Part: s.storage[req.GetUuid()],
	}

	return part, nil
}

func (s *InvService) ListParts(ctx context.Context, req *invServiceV1.ListPartsRequest) (*invServiceV1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filter := req.GetFilter()
	var result []*invServiceV1.Part

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
		if len(filter.ManufacturerCountries) > 0 && !contains(filter.ManufacturerCountries, part.GetManufacturer().GetCountry()) {
			continue
		}

		// Фильтрация по Tags (логическое ИЛИ: хотя бы один тег должен совпасть)
		if len(filter.Tags) > 0 && !hasAnyTag(part.Tags, filter.Tags) {
			continue
		}

		result = append(result, part)
	}

	return &invServiceV1.ListPartsResponse{
		Parts: result,
	}, nil
}

func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func containsCategory(list []invServiceV1.Category, value invServiceV1.Category) bool {
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
