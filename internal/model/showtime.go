package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Showtime struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MovieID     primitive.ObjectID `bson:"movie_id" json:"movie_id"`
	TheaterName string             `bson:"theater_name" json:"theater_name"`
	StartTime   time.Time          `bson:"start_time" json:"start_time"`
	SeatRows    int                `bson:"seat_rows" json:"seat_rows"`
	SeatCols    int                `bson:"seat_cols" json:"seat_cols"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
