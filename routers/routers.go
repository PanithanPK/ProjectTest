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

	accountingRepo := repositorys.NewAccountingRepository(db)
	accountingSvc := services.NewAccountingService(db, userRepo, accountingRepo)
	accountingH := handlers.NewAccountingHandler(accountingSvc)

	RegisterAccountingRoutes(app, auth, accountingH)
}
