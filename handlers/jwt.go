package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func UserIDFromJWT(c *fiber.Ctx) (uint64, error) {
	tokenAny := c.Locals("user")
	token, ok := tokenAny.(*jwt.Token)
	if !ok || token == nil {
		return 0, errors.New("missing jwt token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid jwt claims")
	}
	val, ok := claims["user_id"]
	if !ok {
		return 0, errors.New("user_id claim missing")
	}
	switch v := val.(type) {
	case float64:
		return uint64(v), nil
	case int64:
		return uint64(v), nil
	case int:
		return uint64(v), nil
	case uint64:
		return v, nil
	default:
		return 0, errors.New("invalid user_id claim")
	}
}
