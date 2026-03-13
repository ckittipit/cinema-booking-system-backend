package jobs

import (
	"cinema-booking/backend/internal/service"
	"context"
	"log"
	"time"
)

func StartBookingCleanupJob(bookingService *service.BookingService) {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for range ticker.C {
			if err := bookingService.ExpireTimedOutBookings(context.Background()); err != nil {
				log.Printf("Cleanup job error %v", err)
			}
		}
	}()
}
