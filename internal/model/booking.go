package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingStatus string

const (
	BookingStatusLocked    BookingStatus = "LOCKED"
	BookingStatusConformed BookingStatus = "CONFIRMED"
	BookingStatusExpired   BookingStatus = "EXPIRED"
	BookingStatusCancelled BookingStatus = "CANCELED"
)

type Booking struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	ShowtimeID primitive.ObjectID `bson:"showtime_id" json:"showtime_id"`
	SeatID     string             `bson:"seat_id" json:"seat_id"`
	Status     BookingStatus      `bson:"status" json:"status"`
	LockedAt   *time.Time         `bson:"locked_at,omitempty" json:"locked_at,omitempty"`
	ExpiresAt  *time.Time         `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
	BookedAt   *time.Time         `bson:"booked_at,omitempty" json:"booked_at,omitempty"`
	Price      float64            `bson:"price" json:"price"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}
