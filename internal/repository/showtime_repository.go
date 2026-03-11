package repository

import (
	"cinema-booking/backend/internal/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShowtimeRepository struct {
	collection *mongo.Collection
}

func NewShowtimeRepository(db *mongo.Database) *ShowtimeRepository {
	return &ShowtimeRepository{
		collection: db.Collection("showtimes"),
	}
}

func (r *ShowtimeRepository) FindByMovieID(ctx context.Context, movieID primitive.ObjectID) ([]model.Showtime, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"movie_id": movieID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var showtimes []model.Showtime
	if err := cursor.All(ctx, &showtimes); err != nil {
		return nil, err
	}

	return showtimes, nil
}

func (r *ShowtimeRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Showtime, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var showtime model.Showtime
	if err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&showtime); err != nil {
		return nil, err
	}

	return &showtime, nil
}
