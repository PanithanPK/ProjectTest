package middlewares

import (
	"ProjectTest/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func NewJWTMiddleware(jwtCfg config.JWTConfig) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: jwtCfg.Secret},
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
	})
}
