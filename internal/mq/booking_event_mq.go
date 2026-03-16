package mq

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

const BookingConfirmedChannel = "booking.confirmed"

type BookingConfirmedEvent struct {
	BookingID  string `json:"booking_id"`
	UserID     string `json:"user_id"`
	ShowtimeID string `json:"showtime_id"`
	SeatID     string `json:"seat_id"`
}

type BookingEventMQ struct {
	redisClient *redis.Client
}

func NewBookingEventMQ(redisClient *redis.Client) *BookingEventMQ {
	return &BookingEventMQ{
		redisClient: redisClient,
	}
}

func (m *BookingEventMQ) PublishBookingConfirmed(ctx context.Context, event BookingConfirmedEvent) error {
	raw, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return m.redisClient.Publish(ctx, BookingConfirmedChannel, raw).Err()
}
