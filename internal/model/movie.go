package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Movie struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	DurationMinutes int                `bson:"duration_minutes" json:"duration_minutes"`
	PosterURL       string             `bson:"poster_url" json:"poster_url"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}
