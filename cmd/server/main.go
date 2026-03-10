package main

import (
	"context"
	"log"
	"net/http"

	"cinema-booking/backend/internal/config"
	"cinema-booking/backend/internal/database"
	"cinema-booking/backend/internal/handler"

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

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
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

	healthHandler := handler.NewHealthHandler()
	e.GET("/health", healthHandler.HealthCheck)

	log.Printf("server started on :%s", cfg.AppPort)
	if err := e.Start(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
