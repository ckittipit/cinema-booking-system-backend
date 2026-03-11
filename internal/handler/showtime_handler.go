package handler

import (
	"cinema-booking/backend/internal/service"
	"cinema-booking/backend/internal/utils"

	"github.com/labstack/echo/v4"
)

type ShowtimeHandler struct {
	showtimeService *service.ShowtimeService
}

func NewShowtimeHandler(showtimeService *service.ShowtimeService) *ShowtimeHandler {
	return &ShowtimeHandler{
		showtimeService: showtimeService,
	}
}

func (h *ShowtimeHandler) GetShowtimesByMovieID(c echo.Context) error {
	movieID := c.Param("movieId")
	// println("movieId = ", movieID)

	result, err := h.showtimeService.GetShowtimesByMovieID(c.Request().Context(), movieID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Success(c, "Showtimes fetched successfully", result)
}

func (h *ShowtimeHandler) GetSeatMapByShowtimeID(c echo.Context) error {
	showtimeID := c.Param("showtimeId")

	result, err := h.showtimeService.GetSeatMapByShowtimeID(c.Request().Context(), showtimeID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Success(c, "Seat map fetched successfully", result)
}
