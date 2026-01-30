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

// Login godoc
// @Summary Login user
// @Description User login with username and password
// @Tags User
// @Accept json
// @Produce json
// @Param request body user.LoginRequest true "Login request"
// @Success 200 {object} user.TokenResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /user/login [post]
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

// Register godoc
// @Summary Register user
// @Description Create a new user account
// @Tags User
// @Accept json
// @Produce json
// @Param request body user.RegisterRequest true "Register request"
// @Success 201 {object} user.User
// @Failure 400 {object} utils.ErrorResponse
// @Router /user/register [post]
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

// Me godoc
// @Summary Get current user info
// @Description Get information about the logged-in user
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} user.MeResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /user/me [get]
func (h *UserHandler) Me(c *fiber.Ctx) error {
	userID, err := UserIDFromJWT(c)
	if err != nil {
		return utils.Fail(c, fiber.StatusUnauthorized, "unauthorized")
	}

	res, err := h.svc.Me(c.Context(), userID)
	if err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.Success(c, fiber.StatusOK, res)
}

// Update godoc
// @Summary Update user info
// @Description Update user profile information
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body user.UpdateRequest true "Update request"
// @Success 200 {object} user.User
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /user/update [patch]
func (h *UserHandler) Update(c *fiber.Ctx) error {
	userID, err := UserIDFromJWT(c)
	if err != nil {
		return utils.Fail(c, fiber.StatusUnauthorized, "unauthorized")
	}

	var req user.UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, "invalid request body")
	}

	res, err := h.svc.Update(c.Context(), userID, req)
	if err != nil {
		return utils.Fail(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.Success(c, fiber.StatusOK, res)
}
