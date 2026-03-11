package routes

import (
	"cinema-booking/backend/internal/app"
	"cinema-booking/backend/internal/handler"

	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo, a *app.App) {
	healthHandler := handler.NewHealthHandler()
	movieHandler := handler.NewMovieHandler(a.MovieService)
	showtimeHandler := handler.NewShowtimeHandler(a.ShowtimeService)

	e.GET("/health", healthHandler.HealthCheck)

	api := e.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/movies", movieHandler.GetMovies)
	v1.GET("/movies/:movieId/showtimes", showtimeHandler.GetShowtimesByMovieID)
	v1.GET("/showtimes/:showtimeId/seats", showtimeHandler.GetSeatMapByShowtimeID)
}
