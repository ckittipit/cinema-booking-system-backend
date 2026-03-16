package service

import (
	"cinema-booking/backend/internal/dto"
	"cinema-booking/backend/internal/model"
	"cinema-booking/backend/internal/mq"
	"cinema-booking/backend/internal/repository"
	"cinema-booking/backend/internal/ws"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingService struct {
	bookingRepository  *repository.BookingRepository
	showtimeRepository *repository.ShowtimeRepository
	seatLockService    *SeatLockService
	auditLogService    *AuditLogService
	bookingEventMQ     *mq.BookingEventMQ
	hub                *ws.Hub
}

func NewBookingService(
	bookingRepository *repository.BookingRepository,
	showtimeRepository *repository.ShowtimeRepository,
	seatLockService *SeatLockService,
	auditLogService *AuditLogService,
	hub *ws.Hub,
) *BookingService {
	return &BookingService{
		bookingRepository:  bookingRepository,
		showtimeRepository: showtimeRepository,
		seatLockService:    seatLockService,
		auditLogService:    auditLogService,
		hub:                hub,
	}
}

func (s *BookingService) LockSeat(
	ctx context.Context,
	userID string,
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

	// mockUserID := primitive.NewObjectID()
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	expiresAt := now.Add(SeatLockTTL)

	booking := &model.Booking{
		UserID:     userObjectID,
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

	_ = s.auditLogService.LogEvent(
		ctx,
		"SEAT_LOCKED",
		&booking.UserID,
		&booking.ShowtimeID,
		&booking.SeatID,
		&booking.ID,
		"Seat locked successfully",
		map[string]any{
			"expires_at": expiresAt.Format(time.RFC3339),
		},
	)

	s.hub.BroadcastToShowtime(booking.ShowtimeID.Hex(), map[string]any{
		"type": "seat_locked",
		"data": map[string]any{
			"showtime_id": booking.ShowtimeID.Hex(),
			"seat_id":     booking.SeatID,
		},
	})

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
	userID string,
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

	if booking.UserID.Hex() != userID {
		return nil, fmt.Errorf("forbidden")
	}

	_ = s.bookingEventMQ.PublishBookingConfirmed(ctx, mq.BookingConfirmedEvent{
		BookingID:  booking.ID.Hex(),
		UserID:     booking.UserID.Hex(),
		ShowtimeID: booking.ShowtimeID.Hex(),
		SeatID:     booking.SeatID,
	})

	_ = s.auditLogService.LogEvent(
		ctx,
		"BOOKING_CONFIRMED",
		&booking.UserID,
		&booking.ShowtimeID,
		&booking.SeatID,
		&booking.ID,
		"Booking confirmed successfully",
		nil,
	)

	s.hub.BroadcastToShowtime(booking.ShowtimeID.Hex(), map[string]any{
		"type": "seat_confirmed",
		"data": map[string]any{
			"showtime_id": booking.ShowtimeID.Hex(),
			"seat_id":     booking.SeatID,
		},
	})

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

func (s *BookingService) ReleaseBooking(
	ctx context.Context,
	userID string,
	bookingID string,
) error {
	objectID, err := primitive.ObjectIDFromHex(bookingID)
	if err != nil {
		return fmt.Errorf("Invalid booking id")
	}

	booking, err := s.bookingRepository.FindByID(ctx, objectID)
	if err != nil {
		return fmt.Errorf("Booking not found")
	}

	if booking.Status != model.BookingStatusLocked {
		return fmt.Errorf("Booking is not in locked state")
	}

	if err := s.seatLockService.ReleaseSeatLock(
		ctx,
		booking.ShowtimeID.Hex(),
		booking.SeatID,
	); err != nil {
		return err
	}

	if err := s.bookingRepository.UpdateStatus(ctx, booking.ID, model.BookingStatusCancelled); err != nil {
		return err
	}

	if booking.UserID.Hex() != userID {
		return fmt.Errorf("forbidden")
	}

	_ = s.auditLogService.LogEvent(
		ctx,
		"SEAT_RELEASED",
		&booking.UserID,
		&booking.ShowtimeID,
		&booking.SeatID,
		&booking.ID,
		"Seat released by user",
		nil,
	)

	s.hub.BroadcastToShowtime(booking.ShowtimeID.Hex(), map[string]any{
		"type": "seat_released",
		"data": map[string]any{
			"showtime_id": booking.ShowtimeID.Hex(),
			"seat_id":     booking.SeatID,
		},
	})

	return nil
}

func (s *BookingService) ExpireTimedOutBookings(ctx context.Context) error {
	now := time.Now()

	bookings, err := s.bookingRepository.FindExpiredLockedBookings(ctx, now)
	if err != nil {
		return err
	}

	for _, booking := range bookings {
		_ = s.seatLockService.ReleaseSeatLock(ctx, booking.ShowtimeID.Hex(), booking.SeatID)
		if err := s.bookingRepository.UpdateStatus(ctx, booking.ID, model.BookingStatusExpired); err != nil {
			_ = s.auditLogService.LogEvent(
				ctx,
				"BOOKING_EXPIRED",
				&booking.UserID,
				&booking.ShowtimeID,
				&booking.SeatID,
				&booking.ID,
				"booking expired due to timeout",
				nil,
			)

			s.hub.BroadcastToShowtime(booking.ShowtimeID.Hex(), map[string]any{
				"type": "seat_expired",
				"data": map[string]any{
					"showtime_id": booking.ShowtimeID.Hex(),
					"seat_id":     booking.SeatID,
				},
			})

			return err
		}
	}

	return nil
}
