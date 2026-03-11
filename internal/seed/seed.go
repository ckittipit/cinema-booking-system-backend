package seed

import (
	"cinema-booking/backend/internal/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SeedInitialData(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	movieCollection := db.Collection("movies")
	showtimeCollection := db.Collection("showtimes")

	movieCount, err := movieCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	if movieCount > 0 {
		return nil
	}

	now := time.Now()

	movie1 := model.Movie{
		Title:           "Avengers: Doomsday",
		Description:     "The Avengers must stop Dr.Doom from destroying the multiverse.",
		DurationMinutes: 180,
		PosterURL:       "",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	movie2 := model.Movie{
		Title:           "Spider-Man: Homeless",
		Description:     "Peter Parker is a homeless now, poor Peter",
		DurationMinutes: 120,
		PosterURL:       "",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	movieResult, err := movieCollection.InsertMany(ctx, []any{movie1, movie2})
	if err != nil {
		return err
	}

	movie1ID := movieResult.InsertedIDs[0].(primitive.ObjectID)
	movie2ID := movieResult.InsertedIDs[1].(primitive.ObjectID)

	showtimes := []any{
		model.Showtime{
			MovieID:     movie1ID,
			TheaterName: "Theater 1",
			StartTime:   now.Add(2 * time.Hour),
			SeatRows:    5,
			SeatCols:    8,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		model.Showtime{
			MovieID:     movie1ID,
			TheaterName: "Theater 2",
			StartTime:   now.Add(5 * time.Hour),
			SeatRows:    6,
			SeatCols:    8,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		model.Showtime{
			MovieID:     movie2ID,
			TheaterName: "Theater 4",
			StartTime:   now.Add(1 * time.Hour),
			SeatRows:    10,
			SeatCols:    10,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	_, err = showtimeCollection.InsertMany(ctx, showtimes)
	return err
}
