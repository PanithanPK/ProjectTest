package routers

import (
	"ProjectTest/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app *fiber.App, auth fiber.Handler, h *handlers.UserHandler) {

	app.Post("/user/register", h.Register)
	app.Post("/user/login", h.Login)

	app.Get("/user/me", auth, h.Me)
	app.Patch("/user/update", auth, h.Update)

}
