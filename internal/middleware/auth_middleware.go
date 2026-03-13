package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type CurrentUser struct {
	UserID string
	Role   string
	Email  string
	Name   string
}

func MockAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser := &CurrentUser{
			UserID: "mock-user-id",
			Role:   "USER",
			Email:  "user@example.com",
			Name:   "John Doe",
		}

		c.Set("currentUser", currentUser)
		return next(c)
	}
}

func AdminOnlyModdleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		currentUser, ok := c.Get("currentUser").(*CurrentUser)
		if !ok || currentUser == nil {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"success": false,
				"message": "Unauthorized",
			})
		}

		if currentUser.Role != "ADMIN" {
			return c.JSON(http.StatusForbidden, map[string]any{
				"success": false,
				"message": "Forbidden",
			})
		}

		return next(c)
	}
}
