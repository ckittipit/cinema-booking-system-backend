package app

import (
	"cinema-booking/backend/internal/config"
	"cinema-booking/backend/internal/repository"
	"cinema-booking/backend/internal/service"
	"cinema-booking/backend/internal/ws"
	"log"

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
	UserRepository     *repository.UserRepository
	AuditLogRepository *repository.AuditLogRepository

	MovieService        *service.MovieService
	ShowtimeService     *service.ShowtimeService
	BookingService      *service.BookingService
	SeatLockService     *service.SeatLockService
	AuditLogService     *service.AuditLogService
	AdminService        *service.AdminService
	UserService         *service.UserService
	FirebaseAuthService *service.FirebaseAuthService

	Hub *ws.Hub
}

func New(cfg *config.Config, mongoClient *mongo.Client, redisClient *redis.Client) *App {
	db := mongoClient.Database(cfg.MongoDB)

	movieRepository := repository.NewMovieRepository(db)
	showtimeRepository := repository.NewShowtimeRepository(db)
	bookingRepository := repository.NewBookingRepository(db)
	userRepository := repository.NewUserRepository(db)
	auditLogRepository := repository.NewAuditLogRepository(db)

	seatLockService := service.NewSeatLockService(redisClient)
	auditLogService := service.NewAuditLogService(auditLogRepository)
	userService := service.NewUserService(userRepository)

	firebaseAuthService, err := service.NewFirebaseAuthService(cfg.FirebaseCredentialsPath)
	if err != nil {
		log.Fatalf("failed to init firebase auth service: %v", err)
	}

	hub := ws.NewHub()

	movieService := service.NewMovieService(movieRepository)
	showtimeService := service.NewShowtimeService(showtimeRepository, bookingRepository, seatLockService)
	bookingService := service.NewBookingService(bookingRepository, showtimeRepository, seatLockService, auditLogService, hub)
	adminService := service.NewAdminService(bookingRepository, auditLogService)

	// firebaseAuthService, err := service.NewFirebaseAuthService(cfg.FirebaseCredentialsPath)
	// if err != nil {
	// 	panic(err)
	// }

	return &App{
		Config:              cfg,
		MongoClient:         mongoClient,
		RedisClient:         redisClient,
		Database:            db,
		MovieRepository:     movieRepository,
		ShowtimeRepository:  showtimeRepository,
		BookingRepository:   bookingRepository,
		UserRepository:      userRepository,
		AuditLogRepository:  auditLogRepository,
		MovieService:        movieService,
		ShowtimeService:     showtimeService,
		BookingService:      bookingService,
		SeatLockService:     seatLockService,
		AuditLogService:     auditLogService,
		AdminService:        adminService,
		Hub:                 hub,
		UserService:         userService,
		FirebaseAuthService: firebaseAuthService,
	}
}
