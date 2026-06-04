package middleware

import (
	"fmt"
	"ping-uptime/internal/pkg/jwt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

var jwtService jwt.JWT

func InitializeAuth(service jwt.JWT) {
	jwtService = service
}

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error":   "Authorization header is missing",
				"message": "Unauthorized",
			})
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error":   "Invalid Authorization header format",
				"message": "Unauthorized",
			})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		valid, err := jwtService.ValidateToken(token)
		if err != nil || !valid {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error":   "Invalid or expired token",
				"message": "Unauthorized",
			})
		}

		claims, err := jwtService.ParseToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error":   fmt.Sprintf("Failed to parse token claims: %v", err),
				"message": "Unauthorized",
			})
		}

		c.Set("user", claims)

		return next(c)
	}
}

// Admin middleware enforces that the authenticated user possesses the 'admin' role.
func Admin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(map[string]interface{})
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error":   "Unauthorized",
				"message": "Unauthorized",
			})
		}

		role, ok := user["role"].(string)
		if !ok || role != "admin" {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"error":   "Forbidden: Admin role required",
				"message": "Forbidden",
			})
		}

		return next(c)
	}
}
