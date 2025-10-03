package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
	"github.com/agumiroff/BigTechProject/shared/apperrors"
)

func (r *repository) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	var part model.Part
	err := r.collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	return &part, nil
}
