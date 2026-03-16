package handler

import (
	"cinema-booking/backend/internal/service"
	"cinema-booking/backend/internal/utils"
	"strings"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	firebaseAuthService *service.FirebaseAuthService
	userService         *service.UserService
}

func NewAuthHandler(
	firebaseAuthService *service.FirebaseAuthService,
	userService *service.UserService,
) *AuthHandler {
	return &AuthHandler{
		firebaseAuthService: firebaseAuthService,
		userService:         userService,
	}
}

func (h *AuthHandler) Verify(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return utils.BadRequest(c, "Missing authorization header")
	}

	idToken := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if idToken == authHeader || idToken == "" {
		return utils.BadRequest(c, "Invalid authorization header")
	}

	claims, err := h.firebaseAuthService.VerifyIDToken(c.Request().Context(), idToken)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	uid, _ := claims["user_id"].(string)
	email, _ := claims["email"].(string)
	name, _ := claims["name"].(string)

	if uid == "" {
		return utils.BadRequest(c, "invalid firebase claims")
	}

	user, err := h.userService.FindOrCreateUser(c.Request().Context(), uid, email, name)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.Success(c, "User verified successfully", map[string]any{
		"id":    user.ID.Hex(),
		"uid":   user.FirebaseUID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}
