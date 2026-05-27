package controllers

import "github.com/gofiber/fiber/v2"

func respondError(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"error": message})
}

func respondValidationErrors(c *fiber.Ctx, errors []string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message": "Validation failed",
		"errors":  errors,
	})
}
