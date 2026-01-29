package routers

import (
	"ProjectTest/config"
	"ProjectTest/handlers"
	"ProjectTest/middlewares"
	"ProjectTest/repositorys"
	"ProjectTest/services"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, db *sql.DB, jwtCfg config.JWTConfig) {
	userRepo := repositorys.NewUserRepository(db)
	userSvc := services.NewUserService(userRepo, jwtCfg)
	userH := handlers.NewUserHandler(userSvc)

	auth := middlewares.NewJWTMiddleware(jwtCfg)

	RegisterUserRoutes(app, auth, userH)
}
