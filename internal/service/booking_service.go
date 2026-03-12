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
	seatLockService    *SeatLockService
}

func NewBookingService(
	bookingRepository *repository.BookingRepository,
	showtimeRepository *repository.ShowtimeRepository,
	seatLockService *SeatLockService,
) *BookingService {
	return &BookingService{
		bookingRepository:  bookingRepository,
		showtimeRepository: showtimeRepository,
		seatLockService:    seatLockService,
	}
}

func (s *BookingService) LockSeat(
	ctx context.Context,
	req dto.LockSeatRequest,
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

	existsBooked, err := s.bookingRepository.ExistConfirmedBookingByShowtimeAndSeat(
		ctx,
		showtimeID,
		req.SeatID,
	)

	if err != nil {
		return nil, err
	}
	if existsBooked {
		return nil, fmt.Errorf("Seat booked already")
	}

	now := time.Now()

	existsLocked, err := s.bookingRepository.ExistsActiveLockedBookingByShowTimeAndSeat(
		ctx,
		showtimeID,
		req.SeatID,
		now,
	)
	if err != nil {
		return nil, err
	}
	if existsLocked {
		return nil, fmt.Errorf("Seat locked already")
	}

	mockUserID := primitive.NewObjectID()
	expiresAt := now.Add(SeatLockTTL)

	booking := &model.Booking{
		UserID:     mockUserID,
		ShowtimeID: showtimeID,
		SeatID:     req.SeatID,
		Status:     model.BookingStatusLocked,
		LockedAt:   &now,
		ExpiresAt:  &expiresAt,
		Price:      req.Price,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.bookingRepository.Create(ctx, booking); err != nil {
		return nil, err
	}

	lockValue := SeatLockValue{
		BookingID:  booking.ID.Hex(),
		UserID:     booking.UserID.Hex(),
		ShowtimeID: booking.ShowtimeID.Hex(),
		SeatID:     booking.SeatID,
		ExpiresAt:  expiresAt.Format(time.RFC3339),
	}

	locked, err := s.seatLockService.LockSeat(ctx, req.ShowtimeID, req.SeatID, lockValue)
	if err != nil {
		return nil, err
	}
	if !locked {
		return nil, fmt.Errorf("Seat locked already")
	}

	expiresAtStr := expiresAt.Format("2006-01-02 15:04:05")

	return &dto.BookingResponse{
		ID:         booking.ID.Hex(),
		UserID:     booking.UserID.Hex(),
		ShowtimeID: booking.ShowtimeID.Hex(),
		SeatID:     booking.SeatID,
		Status:     string(booking.Status),
		Price:      booking.Price,
		CreatedAt:  booking.CreatedAt.Format("2006-01-02 15:04:05"),
		ExpiresAt:  &expiresAtStr,
	}, nil
}

func (s *BookingService) ConfirmBooking(
	ctx context.Context,
	bookingID string,
) (*dto.BookingResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(bookingID)
	if err != nil {
		return nil, fmt.Errorf("Invalid booking id")
	}

	booking, err := s.bookingRepository.FindByID(ctx, objectID)
	if err != nil {
		return nil, fmt.Errorf("Booking not found")
	}

	if booking.Status != model.BookingStatusLocked {
		return nil, fmt.Errorf("Booking is not in locked state")
	}

	lockValue, err := s.seatLockService.GetSeatLock(
		ctx,
		booking.ShowtimeID.Hex(),
		booking.SeatID,
	)
	if err != nil {
		return nil, err
	}
	if lockValue == nil {
		return nil, fmt.Errorf("Seat lock expired")
	}
	if lockValue.BookingID != booking.ID.Hex() {
		return nil, fmt.Errorf("Lock does not belong to this booking")
	}

	now := time.Now()
	booking.Status = model.BookingStatusConformed
	booking.BookedAt = &now
	booking.ExpiresAt = nil
	booking.LockedAt = nil

	if err := s.bookingRepository.Update(ctx, booking); err != nil {
		return nil, err
	}

	if err := s.seatLockService.ReleaseSeatLock(
		ctx,
		booking.ShowtimeID.Hex(),
		booking.SeatID,
	); err != nil {
		return nil, err
	}

	return &dto.BookingResponse{
		ID:         booking.ID.Hex(),
		UserID:     booking.UserID.Hex(),
		ShowtimeID: booking.ShowtimeID.Hex(),
		SeatID:     booking.SeatID,
		Status:     string(booking.Status),
		Price:      booking.Price,
		CreatedAt:  booking.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// func (s *BookingService) ConfirmBooking(
// 	ctx context.Context,
// 	req dto.ConfirmBookingRequest,
// ) (*dto.BookingResponse, error) {
// 	showtimeID, err := primitive.ObjectIDFromHex(req.ShowtimeID)
// 	if err != nil {
// 		return nil, fmt.Errorf("Invalid showtime id")
// 	}

// 	showtime, err := s.showtimeRepository.FindByID(ctx, showtimeID)
// 	if err != nil {
// 		return nil, fmt.Errorf("Showtime not found")
// 	}

// 	if !isValidSeatID(req.SeatID, showtime.SeatRows, showtime.SeatCols) {
// 		return nil, fmt.Errorf("Invalid seat id")
// 	}

// 	exists, err := s.bookingRepository.ExistConfirmedBookingByShowtimeAndSeat(
// 		ctx,
// 		showtimeID,
// 		req.SeatID,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if exists {
// 		return nil, fmt.Errorf("Seat already booked")
// 	}

// 	now := time.Now()

// 	mockUserID := primitive.NewObjectID()

// 	booking := &model.Booking{
// 		UserID:     mockUserID,
// 		ShowtimeID: showtimeID,
// 		SeatID:     req.SeatID,
// 		Status:     model.BookingStatusConformed,
// 		BookedAt:   &now,
// 		Price:      req.Price,
// 		CreatedAt:  now,
// 		UpdatedAt:  now,
// 	}

// 	if err := s.bookingRepository.Create(ctx, booking); err != nil {
// 		return nil, err
// 	}

// 	return &dto.BookingResponse{
// 		ID:         booking.ID.Hex(),
// 		UserID:     booking.UserID.Hex(),
// 		ShowtimeID: booking.ShowtimeID.Hex(),
// 		SeatID:     booking.SeatID,
// 		Status:     string(booking.Status),
// 		Price:      booking.Price,
// 		CreatedAt:  booking.CreatedAt.Format("2006-01-02 15-04-05"),
// 	}, nil
// }

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
