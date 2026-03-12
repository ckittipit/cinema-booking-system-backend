package service

import (
	"cinema-booking/backend/internal/dto"
	"cinema-booking/backend/internal/model"
	"cinema-booking/backend/internal/repository"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingService struct {
	bookingRepository  *repository.BookingRepository
	showtimeRepository *repository.ShowtimeRepository
}

func NewBookingService(
	bookingRepository *repository.BookingRepository,
	showtimeRepository *repository.ShowtimeRepository,
) *BookingService {
	return &BookingService{
		bookingRepository:  bookingRepository,
		showtimeRepository: showtimeRepository,
	}
}

func (s *BookingService) ConfirmBooking(
	ctx context.Context,
	req dto.ConfirmBookingRequest,
) (*dto.BookingResponse, error) {
	showtimeID, err := primitive.ObjectIDFromHex(req.ShowtimeID)
	if err != nil {
		return nil, fmt.Errorf("Invalid showtime id")
	}

	showtime, err := s.showtimeRepository.FindByID(ctx, showtimeID)
	if err != nil {
		return nil, fmt.Errorf("Showtime not found")
	}

	if !isValidSeatID(req.SeatID, showtime.SeatRows, showtime.SeatCols) {
		return nil, fmt.Errorf("Invalid seat id")
	}

	exists, err := s.bookingRepository.ExistConfirmedBookingByShowtimeAndSeat(
		ctx,
		showtimeID,
		req.SeatID,
	)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("Seat already booked")
	}

	now := time.Now()

	mockUserID := primitive.NewObjectID()

	booking := &model.Booking{
		UserID:     mockUserID,
		ShowtimeID: showtimeID,
		SeatID:     req.SeatID,
		Status:     model.BookingStatusConformed,
		BookedAt:   &now,
		Price:      req.Price,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.bookingRepository.Create(ctx, booking); err != nil {
		return nil, err
	}

	return &dto.BookingResponse{
		ID:         booking.ID.Hex(),
		UserID:     booking.UserID.Hex(),
		ShowtimeID: booking.ShowtimeID.Hex(),
		SeatID:     booking.SeatID,
		Status:     string(booking.Status),
		Price:      booking.Price,
		CreatedAt:  booking.CreatedAt.Format("2006-01-02 15-04-05"),
	}, nil
}

func isValidSeatID(seatID string, rows int, cols int) bool {
	if len(seatID) < 2 {
		return false
	}

	row := seatID[0]
	if row < 'A' || row >= byte('A'+rows) {
		return false
	}

	var seatNumber int
	_, err := fmt.Sscanf(seatID[1:], "%d", &seatNumber)
	if err != nil {
		return false
	}

	return seatNumber >= 1 && seatNumber <= cols
}
