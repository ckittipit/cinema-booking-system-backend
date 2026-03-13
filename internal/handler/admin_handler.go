package handler

import (
	"cinema-booking/backend/internal/service"
	"cinema-booking/backend/internal/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

func (h *AdminHandler) GetBookings(c echo.Context) error {
	limit := int64(50)

	if raw := c.QueryParam("limit"); raw != "" {
		if parsed, err := strconv.ParseInt(raw, 10, 64); err == nil {
			limit = parsed
		}
	}

	result, err := h.adminService.GetBookings(c.Request().Context(), limit)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.Success(c, "Admin booking fetched successfully", result)
}

func (h *AdminHandler) GetAuditLogs(c echo.Context) error {
	limit := int64(50)

	if raw := c.QueryParam("limit"); raw != "" {
		if parsed, err := strconv.ParseInt(raw, 10, 64); err == nil {
			limit = parsed
		}
	}

	result, err := h.adminService.GetAuditLogs(c.Request().Context(), limit)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.Success(c, "Admin audit logs fetched successfully", result)
}
