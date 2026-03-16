package main

import (
	"context"
	"log"
	"net/http"

	"cinema-booking/backend/internal/app"
	"cinema-booking/backend/internal/config"
	"cinema-booking/backend/internal/database"
	"cinema-booking/backend/internal/jobs"
	"cinema-booking/backend/internal/routes"
	"cinema-booking/backend/internal/seed"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	mongoClient, err := database.NewMongoClient(cfg)
	if err != nil {
		log.Fatalf("failed to connect mongo: %v", err)
	}
	defer func() {
		_ = mongoClient.Disconnect(context.Background())
	}()

	redisClient, err := database.NewRedisClient(cfg)
	if err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}
	defer func() {
		_ = redisClient.Close()
	}()

	a := app.New(cfg, mongoClient, redisClient)
	jobs.StartBookingEventConsumer(a.RedisClient, a.AuditLogService)
	jobs.StartBookingCleanupJob(a.BookingService)

	if err := seed.SeedInitialData(a.Database); err != nil {
		log.Fatalf("failed to seed data: %v", err)
	}

	e := echo.New()
	e.HideBanner = true

	// e.Use(middleware.Logger())
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{cfg.FrontendOrigin},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))

	routes.Register(e, a)
	// healthHandler := handler.NewHealthHandler()
	// e.GET("/health", healthHandler.HealthCheck)

	log.Printf("server started on :%s", cfg.AppPort)
	if err := e.Start(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
