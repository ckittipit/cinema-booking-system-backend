package service

import (
	"cinema-booking/backend/internal/model"
	"cinema-booking/backend/internal/repository"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuditLogService struct {
	auditLogRepository *repository.AuditLogRepository
}

func NewAuditLogService(auditLogRepository *repository.AuditLogRepository) *AuditLogService {
	return &AuditLogService{
		auditLogRepository: auditLogRepository,
	}
}

func (s *AuditLogService) LogEvent(
	ctx context.Context,
	eventType string,
	userID *primitive.ObjectID,
	showtimeID *primitive.ObjectID,
	seatID *string,
	bookingID *primitive.ObjectID,
	message string,
	metadata map[string]any,
) error {
	log := &model.AuditLog{
		EventType:  eventType,
		UserID:     userID,
		ShowtimeID: showtimeID,
		SeatID:     *seatID,
		BookingID:  bookingID,
		Message:    message,
		Metadata:   metadata,
		CreatedAt:  time.Now(),
	}

	return s.auditLogRepository.Create(ctx, log)
}

func (s *AuditLogService) GetAuditLogs(ctx context.Context, limit int64) ([]model.AuditLog, error) {
	return s.auditLogRepository.FindAll(ctx, limit)
}
