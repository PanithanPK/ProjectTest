package routers

import (
	"ProjectTest/handlers"
	"ProjectTest/repositorys"
	"ProjectTest/services"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, db *sql.DB) {
	userRepo := repositorys.NewUserRepository(db)
	userSvc := services.NewUserService(userRepo)
	userH := handlers.NewUserHandler(userSvc)

	RegisterUserRoutes(app, userH)
}
