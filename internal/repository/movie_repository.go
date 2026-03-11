package repository

import (
	"cinema-booking/backend/internal/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MovieRepository struct {
	collection *mongo.Collection
}

func NewMovieRepository(db *mongo.Database) *MovieRepository {
	return &MovieRepository{
		collection: db.Collection("movies"),
	}
}

func (r *MovieRepository) FindAll(ctx context.Context) ([]model.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var movies []model.Movie
	if err := cursor.All(ctx, &movies); err != nil {
		return nil, err
	}

	return movies, nil
}
