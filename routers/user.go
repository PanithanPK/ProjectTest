package routers

import (
	"ProjectTest/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app *fiber.App, auth fiber.Handler, h *handlers.UserHandler) {

	app.Post("/user/login", h.Login)

}
