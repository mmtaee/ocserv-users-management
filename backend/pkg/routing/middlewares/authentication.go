package middlewares

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"ocserv-bakend/pkg/config"
	"strings"
)

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cfg := config.Get()
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return UnauthorizedError(c, "missing or invalid Authorization header")
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(cfg.JWTSecret), nil
			})

			if err != nil || !token.Valid {
				return UnauthorizedError(c, "invalid or expired token")
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if userUID, ok := claims["sub"].(string); ok {
					c.Set("userUID", userUID)
					c.Set("isAdmin", claims["isAdmin"])
				} else {
					return UnauthorizedError(c, "user ID not found in token")
				}
			} else {
				return UnauthorizedError(c, "invalid token claims")
			}

			return next(c)
		}
	}
}
