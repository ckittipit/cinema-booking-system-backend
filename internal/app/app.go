package app

import (
	"cinema-booking/backend/internal/config"
	"cinema-booking/backend/internal/repository"
	"cinema-booking/backend/internal/service"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Config *config.Config

	MongoClient *mongo.Client
	RedisClient *redis.Client
	Database    *mongo.Database

	MovieRepository    *repository.MovieRepository
	ShowtimeRepository *repository.ShowtimeRepository
	BookingRepository  *repository.BookingRepository

	MovieService    *service.MovieService
	ShowtimeService *service.ShowtimeService
	BookingService  *service.BookingService
	SeatLockService *service.SeatLockService
}

func New(cfg *config.Config, mongoClient *mongo.Client, redisClient *redis.Client) *App {
	db := mongoClient.Database(cfg.MongoDB)

	movieRepository := repository.NewMovieRepository(db)
	showtimeRepository := repository.NewShowtimeRepository(db)
	bookingRepository := repository.NewBookingRepository(db)

	seatLockService := service.NewSeatLockService(redisClient)
	movieService := service.NewMovieService(movieRepository)
	showtimeService := service.NewShowtimeService(showtimeRepository, bookingRepository, seatLockService)
	bookingService := service.NewBookingService(bookingRepository, showtimeRepository, seatLockService)

	return &App{
		Config:             cfg,
		MongoClient:        mongoClient,
		RedisClient:        redisClient,
		Database:           db,
		MovieRepository:    movieRepository,
		ShowtimeRepository: showtimeRepository,
		BookingRepository:  bookingRepository,
		MovieService:       movieService,
		ShowtimeService:    showtimeService,
		BookingService:     bookingService,
	}
}
