package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
	"github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/converter"
	repomodel "github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/model"
)

func (s *repository) ListParts(ctx context.Context, filter *repomodel.PartsFilter) ([]*model.Part, error) {
	query := bson.M{}

	if len(filter.UUIDs) > 0 {
		query["uuid"] = bson.M{"$in": filter.UUIDs}
	}

	if len(filter.Names) > 0 {
		query["name"] = bson.M{"$in": filter.Names}
	}

	if len(filter.Categories) > 0 {
		query["category"] = bson.M{"$in": filter.Categories}
	}

	if len(filter.ManufacturerCountries) > 0 {
		query["manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
	}

	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$in": filter.Tags}
	}

	cursor, err := s.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var repoParts []*repomodel.Part
	if err = cursor.All(ctx, &repoParts); err != nil {
		return nil, err
	}

	parts := make([]*model.Part, len(repoParts))
	for i, part := range repoParts {
		parts[i] = converter.RepoToModel(part)
	}

	return parts, nil
}
