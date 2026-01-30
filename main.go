package main

import (
	"ProjectTest/config"
	"ProjectTest/routers"
	"ProjectTest/utils"
	"database/sql"
	"log"
	"os"

	_ "ProjectTest/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title ProjectTest API
// @version 1.0
// @description Money Transfer API
// @host localhost:5000
// @basePath /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	_ = godotenv.Load()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "5000"
	}

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	jwtCfg, err := config.LoadJWTConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()
	app.Use(recover.New())

	// Swagger route
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	routers.Register(app, db, jwtCfg)

	app.Use(func(c *fiber.Ctx) error {
		return utils.Fail(c, fiber.StatusNotFound, "route not found")
	})

	log.Fatal(app.Listen(":" + port))
}
