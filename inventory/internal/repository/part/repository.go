package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	rep "github.com/agumiroff/BigTechProject/inventory/v1/internal/repository"
)

type repository struct {
	collection *mongo.Collection
}

var _ rep.InvRepository = (*repository)(nil)

func NewRepository(db *mongo.Database) *repository {
	collection := db.Collection("inventory")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "title", Value: 1}},
			Options: options.Index().SetUnique(false),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		panic(err)
	}

	return &repository{
		collection: collection,
	}
}
