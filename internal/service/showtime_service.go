package service

import (
	"cinema-booking/backend/internal/dto"
	"cinema-booking/backend/internal/repository"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShowtimeService struct {
	shotimeRepository *repository.ShowtimeRepository
	bookingRepository *repository.BookingRepository
	seatLockService   *SeatLockService
}

func NewShowtimeService(
	showtimeRepository *repository.ShowtimeRepository,
	bookingRepository *repository.BookingRepository,
	seatLockService *SeatLockService,
) *ShowtimeService {
	return &ShowtimeService{
		shotimeRepository: showtimeRepository,
		bookingRepository: bookingRepository,
		seatLockService:   seatLockService,
	}
}

// func (s *ShowtimeService) GetShowtimesByMovieID(ctx context.Context, movieID string) (any, error) {
// 	objectID, err := primitive.ObjectIDFromHex(movieID)
// 	if err != nil {
// 		return nil, fmt.Errorf("Invalid movie id")
// 	}

//		return s.shotimeRepository.FindByMovieID(ctx, objectID)
//	}
func (s *ShowtimeService) GetShowtimesByMovieID(ctx context.Context, movieID string) ([]dto.ShowtimeResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		return nil, fmt.Errorf("Invalid movie id")
	}

	showtimes, err := s.shotimeRepository.FindByMovieID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.ShowtimeResponse, 0, len(showtimes))
	for _, showtime := range showtimes {
		result = append(result, dto.ShowtimeResponse{
			ID:          showtime.ID.Hex(),
			MovieID:     showtime.MovieID.Hex(),
			TheaterName: showtime.TheaterName,
			StartTime:   showtime.StartTime.Format("2006-01-02 15:04:05"),
			SeatRows:    showtime.SeatRows,
			SeatCols:    showtime.SeatCols,
		})
	}

	return result, nil
}

func (s *ShowtimeService) GetSeatMapByShowtimeID(ctx context.Context, showtimeID string) (*dto.SeatMapResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(showtimeID)
	if err != nil {
		return nil, fmt.Errorf("Invalid showtime id")
	}

	showtime, err := s.shotimeRepository.FindByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	bookedSeatIDs, err := s.bookingRepository.FindConfirmSeatIDsByShowtimeID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	lockedBookings, err := s.bookingRepository.FindActiveLockedBookingsByShowtimeID(ctx, objectID, time.Now())
	if err != nil {
		return nil, err
	}
	// to debug
	// fmt.Println("Locked bookings count =", len(lockedBookings))
	// for _, booking := range lockedBookings {
	// 	fmt.Println("Locked seat = ", booking.SeatID, "Status =", booking.Status, "ExpiresAt =", booking.ExpiresAt)
	// }

	bookedSet := make(map[string]bool, len(bookedSeatIDs))
	for _, seatID := range bookedSeatIDs {
		bookedSet[seatID] = true
	}

	lockedExpiryMap := make(map[string]string, len(lockedBookings))
	for _, booking := range lockedBookings {
		if booking.ExpiresAt != nil {
			lockedExpiryMap[booking.SeatID] = booking.ExpiresAt.Format("2006-01-02 15:04:05")
		}
	}

	seats := make([]dto.SeatResponse, 0)

	for row := 0; row < showtime.SeatRows; row++ {
		rowLetter := string(rune('A' + row))
		for col := 1; col <= showtime.SeatCols; col++ {
			seatID := fmt.Sprintf("%s%d", rowLetter, col)
			status := "AVAILABLE"
			var expiresAt *string

			if bookedSet[seatID] {
				status = "BOOKED"
			} else if expiry, ok := lockedExpiryMap[seatID]; ok {
				status = "LOCKED"
				expiresAt = &expiry
			}

			seats = append(seats, dto.SeatResponse{
				SeatID:    seatID,
				Status:    status,
				ExpiresAt: expiresAt,
			})
		}
	}

	return &dto.SeatMapResponse{
		ShowtimeID: showtime.ID.Hex(),
		Seats:      seats,
	}, nil
}
