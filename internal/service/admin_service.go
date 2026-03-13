package service

import (
	"cinema-booking/backend/internal/model"
	"cinema-booking/backend/internal/repository"
	"context"
)

type AdminService struct {
	bookingRepository *repository.BookingRepository
	AuditLogService   *AuditLogService
}

func NewAdminService(
	bookingRepository *repository.BookingRepository,
	auditLogService *AuditLogService,
) *AdminService {
	return &AdminService{
		bookingRepository: bookingRepository,
		AuditLogService:   auditLogService,
	}
}

func (s *AdminService) GetBookings(ctx context.Context, limit int64) ([]model.Booking, error) {
	return s.bookingRepository.FindAll(ctx, limit)
}

func (s *AdminService) GetAuditLogs(ctx context.Context, limit int64) ([]model.AuditLog, error) {
	return s.AuditLogService.GetAuditLogs(ctx, limit)
}
