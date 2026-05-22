package controllers

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"go-fiber-test/database"
	m "go-fiber-test/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func GetProfileUserEndpoint(c *fiber.Ctx) error {
	var profileUser []m.ProfileUser

	if err := database.DBConn.Find(&profileUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch profile user"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"count": len(profileUser),
		"data":  profileUser,
	})

}

func AddProfileUserEndpoint(c *fiber.Ctx) error {
	var profileUser m.ProfileUser

	if err := c.BodyParser(&profileUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := validate.Struct(&profileUser); err != nil {
		errors := formatValidationErrors(err.(validator.ValidationErrors))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	if err := database.DBConn.Create(&profileUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(profileUser)

}

func UpdateProfileUserEndpoint(c *fiber.Ctx) error {
	id := c.Params("id")

	var profileUser m.ProfileUser
	if err := c.BodyParser(&profileUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error" : "Cannot parse JSON"})
	}

	if err := validate.Struct(&profileUser); err != nil {
		errors := formatValidationErrors(err.(validator.ValidationErrors))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	result := database.DBConn.Model(&m.ProfileUser{}).Where("id = ?", id).Updates(&profileUser)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error" : "Could not update profile user"})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error" : "Profile user not found"})
	}

	return c.Status(fiber.StatusOK).JSON(profileUser)

}

func RemoveProfileUserEndpoint(c *fiber.Ctx) error {
	id := c.Params("id")

	result := database.DBConn.Delete(&m.ProfileUser{}, id)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error" : "Could not delete profile user"})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error" : "Profile user not found"})
	}

	return c.SendStatus(fiber.StatusNoContent)

}

func GetUserGenerationsEndpoint(c *fiber.Ctx) error {

	var users []m.ProfileUser
	
	if err := database.DBConn.Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch users"})
	}

	result := m.UserGenerationResult{
		Data:  users,
		Count: len(users),
	}

	for _, u := range users {
		if u.Age < 24 {
			result.GenZ++
		} else if u.Age <= 41 {
			result.GenY++
		} else if u.Age <= 56 {
			result.GenX++
		} else if u.Age <= 75 {
			result.BabyBoomer++
		} else {
			result.GIGeneration++
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func SearchProfileUsersEndpoint(c *fiber.Ctx) error {
	
	search := strings.TrimSpace(c.Query("search"))

	var users []m.ProfileUser

	if err := database.DBConn.Scopes(m.SearchUserScope(search)).Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not search profile users"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"count": len(users),
		"data":  users,
	})
}

func SeedProfileUsersEndpoint(c *fiber.Ctx) error {
	if err := database.DBConn.Unscoped().Where("1 = 1").Delete(&m.ProfileUser{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not clear users"})
	}

	rand.Seed(time.Now().UnixNano())
	
	testUsers := make([]m.ProfileUser, 0, 25)
	for i := 1; i <= 25; i++ {
		age := rand.Intn(81) + 10 // สุ่มอายุระหว่าง 10 - 90 ปี
		testUsers = append(testUsers, m.ProfileUser{
			EmployeeID: fmt.Sprintf("EMP-%03d", i),
			Name:       fmt.Sprintf("User-%d", i),
			Lastname:   fmt.Sprintf("Lastname-%d", i),
			Birthday:   "1990-01-01",
			Age:        age,
			Email:      fmt.Sprintf("user%d@example.com", i),
			Telephone:  fmt.Sprintf("0810000%03d", i),
		})
	}

	if err := database.DBConn.Create(&testUsers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Seed data created successfully", "count": len(testUsers)})
}