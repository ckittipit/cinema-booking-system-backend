package service

import (
	"cinema-booking/backend/internal/dto"
	"cinema-booking/backend/internal/repository"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShowtimeService struct {
	shotimeRepository *repository.ShowtimeRepository
	bookingRepository *repository.BookingRepository
}

func NewShowtimeService(showtimeRepository *repository.ShowtimeRepository, bookingRepository *repository.BookingRepository) *ShowtimeService {
	return &ShowtimeService{
		shotimeRepository: showtimeRepository,
		bookingRepository: bookingRepository,
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

	bookedSet := make(map[string]bool, len(bookedSeatIDs))
	for _, seatID := range bookedSeatIDs {
		bookedSet[seatID] = true
	}

	seats := make([]dto.SeatResponse, 0)

	for row := 0; row < showtime.SeatRows; row++ {
		rowLetter := string(rune('A' + row))
		for col := 1; col <= showtime.SeatCols; col++ {
			seatID := fmt.Sprintf("%s%d", rowLetter, col)
			status := "AVAILABLE"

			if bookedSet[seatID] {
				status = "BOOKED"
			}

			seats = append(seats, dto.SeatResponse{
				SeatID: seatID,
				Status: status,
			})
		}
	}

	return &dto.SeatMapResponse{
		ShowtimeID: showtime.ID.Hex(),
		Seats:      seats,
	}, nil
}
