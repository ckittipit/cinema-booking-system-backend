package handler

import (
	"cinema-booking/backend/internal/dto"
	"cinema-booking/backend/internal/service"
	"cinema-booking/backend/internal/utils"

	"github.com/labstack/echo/v4"
)

type BookingHandler struct {
	bookingService *service.BookingService
}

func NewBookingHandler(bookingService *service.BookingService) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
	}
}

func (h *BookingHandler) LockSeat(c echo.Context) error {
	var req dto.LockSeatRequest

	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	if req.ShowtimeID == "" || req.SeatID == "" {
		return utils.BadRequest(c, "showtime_id and seat_id are required")
	}

	result, err := h.bookingService.LockSeat(c.Request().Context(), req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Created(c, "Seat locked successfully", result)
}

func (h *BookingHandler) ConfirmBooking(c echo.Context) error {
	bookingID := c.Param("bookingId")
	if bookingID == "" {
		return utils.BadRequest(c, "bookingId is required")
	}

	result, err := h.bookingService.ConfirmBooking(c.Request().Context(), bookingID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Success(c, "Booking confirmed successfully", result)
}

// func (h *BookingHandler) ConfirmBooking(c echo.Context) error {
// 	var req dto.ConfirmBookingRequest

// 	if err := c.Bind(&req); err != nil {
// 		return utils.BadRequest(c, "Invalid request body")
// 	}

// 	if req.ShowtimeID == "" || req.SeatID == "" {
// 		return utils.BadRequest(c, "showtime_id and seat_id are required")
// 	}

// 	result, err := h.bookingService.ConfirmBooking(c.Request().Context(), req)
// 	if err != nil {
// 		return utils.BadRequest(c, err.Error())
// 	}

// 	return utils.Created(c, "Booking confirmed successfully", result)
// }

func (h *BookingHandler) ReleaseBooking(c echo.Context) error {
	bookingID := c.Param("bookingId")
	if bookingID == "" {
		return utils.BadRequest(c, "bookingId is required")
	}

	if err := h.bookingService.ReleaseBooking(c.Request().Context(), bookingID); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Success(c, "Booking release successfully", nil)
}
