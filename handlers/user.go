package handlers

import (
	"ProjectTest/modules/user"
	"ProjectTest/services"
	"ProjectTest/utils"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	svc *services.UserService
}

func NewUserHandler(svc *services.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req user.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalid request body")
	}

	res, err := h.svc.Login(c.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			return utils.Fail(c, fiber.StatusUnauthorized, err.Error())
		}
		return utils.Fail(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.Success(c, fiber.StatusOK, res)
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req user.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalid request body")
	}

	res, err := h.svc.Register(c.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUsernameTaken), errors.Is(err, services.ErrBankAccountTaken), errors.Is(err, services.ErrInvalidBankAccount):
			return utils.Fail(c, fiber.StatusBadRequest, err.Error())
		default:
			return utils.Fail(c, fiber.StatusBadRequest, err.Error())
		}
	}
	return utils.Success(c, fiber.StatusCreated, res)
}
