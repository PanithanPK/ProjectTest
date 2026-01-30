package handlers

import (
	"ProjectTest/modules/accounting"
	"ProjectTest/services"
	"ProjectTest/utils"

	"github.com/gofiber/fiber/v2"
)

type AccountingHandler struct {
	svc *services.AccountingService
}

func NewAccountingHandler(svc *services.AccountingService) *AccountingHandler {
	return &AccountingHandler{svc: svc}
}

// Transfer godoc
// @Summary Create transfer
// @Description Transfer money to another user
// @Tags Accounting
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body accounting.TransferRequest true "Transfer request"
// @Success 200 {object} map[string]int64
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /accounting/transfer [post]
func (h *AccountingHandler) Transfer(c *fiber.Ctx) error {
	userID, err := UserIDFromJWT(c)
	if err != nil {
		return utils.Fail(c, fiber.StatusUnauthorized, "unauthorized")
	}

	var req accounting.TransferRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalid request body")
	}

	transferID, err := h.svc.Transfer(c.Context(), userID, req)
	if err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.Success(c, fiber.StatusOK, fiber.Map{"transfer_id": transferID})
}

// TransferList godoc
// @Summary List transfers
// @Description Get list of transfers for the current user
// @Tags Accounting
// @Accept json
// @Produce json
// @Security Bearer
// @Param start_date query string false "Start date"
// @Param end_date query string false "End date"
// @Success 200 {array} accounting.TransferItem
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /accounting/transfer-list [get]
func (h *AccountingHandler) TransferList(c *fiber.Ctx) error {
	userID, err := UserIDFromJWT(c)
	if err != nil {
		return utils.Fail(c, fiber.StatusUnauthorized, "unauthorized")
	}

	start := c.Query("start_date")
	end := c.Query("end_date")

	items, err := h.svc.TransferList(c.Context(), userID, start, end)
	if err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.Success(c, fiber.StatusOK, items)
}
