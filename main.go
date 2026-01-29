package main

import (
	"ProjectTest/router"
	"ProjectTest/utils"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "5000"
	}

	app := fiber.New()
	app.Use(recover.New())

	router.Register(app)

	app.Use(func(c *fiber.Ctx) error {
		return utils.Fail(c, fiber.StatusNotFound, "route not found")
	})

	log.Fatal(app.Listen(":" + port))
}
