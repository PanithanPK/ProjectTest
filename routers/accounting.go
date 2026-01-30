package routers

import (
	"ProjectTest/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterAccountingRoutes(app *fiber.App, auth fiber.Handler, h *handlers.AccountingHandler) {
	app.Post("/accounting/transfer", auth, h.Transfer)
	app.Get("/accounting/transfer-list", auth, h.TransferList)
}
