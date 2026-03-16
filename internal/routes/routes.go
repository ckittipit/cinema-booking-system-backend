package routes

import (
	"cinema-booking/backend/internal/app"
	"cinema-booking/backend/internal/handler"

	"github.com/labstack/echo/v4"

	appmw "cinema-booking/backend/internal/middleware"
)

func Register(e *echo.Echo, a *app.App) {
	healthHandler := handler.NewHealthHandler()
	movieHandler := handler.NewMovieHandler(a.MovieService)
	showtimeHandler := handler.NewShowtimeHandler(a.ShowtimeService)
	bookingHandler := handler.NewBookingHandler(a.BookingService)
	wsHandler := handler.NewWSHandler(a.Hub)
	adminHandler := handler.NewAdminHandler(a.AdminService)
	authHandler := handler.NewAuthHandler(a.FirebaseAuthService, a.UserService)

	authMiddleware := appmw.NewAuthMiddleware(a.FirebaseAuthService, a.UserService)

	e.GET("/health", healthHandler.HealthCheck)
	e.GET("/ws", wsHandler.Handle)

	api := e.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/movies", movieHandler.GetMovies)
	v1.GET("/movies/:movieId/showtimes", showtimeHandler.GetShowtimesByMovieID)
	v1.GET("/showtimes/:showtimeId/seats", showtimeHandler.GetSeatMapByShowtimeID)

	v1.POST("/auth/verify", authHandler.Verify)

	bookings := v1.Group("/bookings")
	bookings.Use(authMiddleware.RequireAuth)
	v1.POST("/bookings/lock", bookingHandler.LockSeat)
	v1.POST("/bookings/:bookingId/confirm", bookingHandler.ConfirmBooking)
	v1.POST("/bookings/:bookingId/release", bookingHandler.ReleaseBooking)

	admin := v1.Group("/admin")
	admin.Use(authMiddleware.RequireAuth)
	// admin.Use(appmw.MockAuthMiddleware)
	admin.Use(appmw.AdminOnlyModdleware)
	admin.GET("/bookings", adminHandler.GetBookings)
	admin.GET("/audit-logs", adminHandler.GetAuditLogs)
}
