package router

import "github.com/gofiber/fiber/v2"

func Register(app *fiber.App) {
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("hello")
	})
}
