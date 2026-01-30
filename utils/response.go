package utils

import "github.com/gofiber/fiber/v2"

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func Success(c *fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(APIResponse{Success: true, Data: data})
}

func Fail(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(APIResponse{Success: false, Message: message})
}
