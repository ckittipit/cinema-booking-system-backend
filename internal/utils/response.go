package utils

import (
	"cinema-booking/backend/internal/dto"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Success(c echo.Context, message string, data any) error {
	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Created(c echo.Context, message string, data any) error {
	return c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func BadRequest(c echo.Context, message string) error {
	return c.JSON(http.StatusBadRequest, dto.APIResponse{
		Success: false,
		Message: message,
	})
}

func InternalError(c echo.Context, message string) error {
	return c.JSON(http.StatusInternalServerError, dto.APIResponse{
		Success: false,
		Message: message,
	})
}

func NotFound(e echo.Context, message string) error {
	return e.JSON(http.StatusNotFound, dto.APIResponse{
		Success: false,
		Message: message,
	})
}
