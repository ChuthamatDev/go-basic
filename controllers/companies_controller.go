package controllers

import (
	"go-fiber-test/database"
	m "go-fiber-test/models"

	"github.com/gofiber/fiber/v2"
)

func GetCompaniesEndpoint(c *fiber.Ctx) error {
	db := database.DBConn
	var companies []m.Company
	db.Find(&companies)
	return c.Status(fiber.StatusOK).JSON(companies)
}

func GetCompanyByIdEndpoint(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")
	var company m.Company
	result := db.First(&company, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Company not found"})
	}
	return c.Status(fiber.StatusOK).JSON(company)
}

func AddCompanyEndpoint(c *fiber.Ctx) error {
	db := database.DBConn
	var company m.Company
	if err := c.BodyParser(&company); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	if err := db.Create(&company).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create company"})
	}
	return c.Status(fiber.StatusCreated).JSON(company)
}

func UpdateCompanyEndpoint(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")
	var company m.Company
	if err := c.BodyParser(&company); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	if err := db.Where("id = ?", id).Updates(&company).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update company"})
	}
	return c.Status(fiber.StatusOK).JSON(company)
}

func RemoveCompanyEndpoint(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")
	var company m.Company
	result := db.Delete(&company, id)
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Company not found"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}