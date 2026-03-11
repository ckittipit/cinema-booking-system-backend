package handler

import (
	"cinema-booking/backend/internal/service"
	"cinema-booking/backend/internal/utils"

	"github.com/labstack/echo/v4"
)

type MovieHandler struct {
	movieService *service.MovieService
}

func NewMovieHandler(movieService *service.MovieService) *MovieHandler {
	return &MovieHandler{
		movieService: movieService,
	}
}

func (h *MovieHandler) GetMovies(c echo.Context) error {
	result, err := h.movieService.GetMovies(c.Request().Context())
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.Success(c, "Movies fetched successfully", result)
}
