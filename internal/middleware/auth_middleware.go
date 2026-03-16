package middleware

import (
	"cinema-booking/backend/internal/service"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type CurrentUser struct {
	UserID string
	Role   string
	Email  string
	Name   string
}

type AuthMiddleware struct {
	firebaseAuthService *service.FirebaseAuthService
	userService         *service.UserService
}

func NewAuthMiddleware(
	firebaseAuthService *service.FirebaseAuthService,
	userService *service.UserService,
) *AuthMiddleware {
	return &AuthMiddleware{
		firebaseAuthService: firebaseAuthService,
		userService:         userService,
	}
}

// func MockAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		currentUser := &CurrentUser{
// 			UserID: "mock-user-id",
// 			Role:   "ADMIN",
// 			Email:  "user@example.com",
// 			Name:   "John Doe",
// 		}

// 		c.Set("currentUser", currentUser)
// 		return next(c)
// 	}
// }

func (m *AuthMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"success": false,
				"message": "Missingauthorization header",
			})
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		// fmt.Println("auth header exists:", authHeader != "")
		// fmt.Println("token prefix ok:", idToken != authHeader)
		// fmt.Println("token dot count:", strings.Count(idToken, "."))
		if idToken == authHeader || strings.TrimSpace(idToken) == "" {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"success": false,
				"message": "Invalid authorization header",
			})
		}

		claims, err := m.firebaseAuthService.VerifyIDToken(c.Request().Context(), idToken)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"success": false,
				"message": "Invalid firebase token",
			})
		}

		uid, _ := claims["user_id"].(string)
		email, _ := claims["email"].(string)
		name, _ := claims["name"].(string)

		if uid == "" {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"success": false,
				"message": "Invalid firebase claims",
			})
		}

		user, err := m.userService.FindOrCreateUser(c.Request().Context(), uid, email, name)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"success": false,
				"message": "Failed to load user",
			})
		}

		currentUser := &CurrentUser{
			UserID: user.ID.Hex(),
			Role:   string(user.Role),
			Email:  user.Email,
			Name:   user.Name,
		}

		c.Set("currentUser", currentUser)
		return next(c)
	}
}

func (m *AuthMiddleware) OptionalAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return next(c)
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		if idToken == authHeader || strings.TrimSpace(idToken) == "" {
			return next(c)
		}

		claims, err := m.firebaseAuthService.VerifyIDToken(c.Request().Context(), idToken)
		if err != nil {
			return next(c)
		}

		uid, _ := claims["user_id"].(string)
		email, _ := claims["email"].(string)
		name, _ := claims["name"].(string)

		if uid == "" {
			return next(c)
		}

		user, err := m.userService.FindOrCreateUser(c.Request().Context(), uid, email, name)
		if err != nil {
			return next(c)
		}

		currentUser := &CurrentUser{
			UserID: user.ID.Hex(),
			Role:   string(user.Role),
			Email:  user.Email,
			Name:   user.Name,
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
