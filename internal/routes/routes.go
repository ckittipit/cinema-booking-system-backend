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
	bookingHandler := handler.NewBookingHandler(a.BookingService)

	e.GET("/health", healthHandler.HealthCheck)

	api := e.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/movies", movieHandler.GetMovies)
	v1.GET("/movies/:movieId/showtimes", showtimeHandler.GetShowtimesByMovieID)
	v1.GET("/showtimes/:showtimeId/seats", showtimeHandler.GetSeatMapByShowtimeID)

	v1.POST("/bookings/lock", bookingHandler.LockSeat)
	v1.POST("/bookings/:bookingId/confirm", bookingHandler.ConfirmBooking)
	v1.POST("bookings/:bookingId/release", bookingHandler.ReleaseBooking)
}
