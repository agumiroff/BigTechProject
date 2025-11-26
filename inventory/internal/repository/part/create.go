package repository

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
)

func (r *repository) CreatePart(ctx context.Context, part *model.Part) (string, error) {
	part.Uuid = gofakeit.UUID()
	part.CreatedAt = time.Now().Unix()
	part.UpdatedAt = part.CreatedAt
	_, err := r.collection.InsertOne(ctx, part)
	if err != nil {
		return "", err
	}

	return part.Uuid, nil
}
