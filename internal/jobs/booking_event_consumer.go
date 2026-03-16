package jobs

import (
	"cinema-booking/backend/internal/mq"
	"cinema-booking/backend/internal/service"
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StartBookingEventConsumer(
	redisClient *redis.Client,
	auditLogService *service.AuditLogService,
) {
	pubsub := redisClient.Subscribe(context.Background(), mq.BookingConfirmedChannel)

	go func() {
		ch := pubsub.Channel()

		for msg := range ch {
			var event mq.BookingConfirmedEvent
			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				log.Printf("Failed to unmarshal booking event: %v", err)
				continue
			}

			var userIDPtr *primitive.ObjectID
			var showtimeIDPtr *primitive.ObjectID
			var bookingIDPtr *primitive.ObjectID

			if id, err := primitive.ObjectIDFromHex(event.UserID); err == nil {
				userIDPtr = &id
			}
			if id, err := primitive.ObjectIDFromHex(event.ShowtimeID); err == nil {
				showtimeIDPtr = &id
			}
			if id, err := primitive.ObjectIDFromHex(event.BookingID); err == nil {
				bookingIDPtr = &id
			}

			_ = auditLogService.LogEvent(
				context.Background(),
				"BOOKING_CONFIRMED_ASYNC",
				userIDPtr,
				showtimeIDPtr,
				event.SeatID,
				bookingIDPtr,
				"Async consumer processed booking confirmed event",
				nil,
			)
		}
	}()
}
