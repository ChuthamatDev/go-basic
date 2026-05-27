package controllers

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"go-fiber-test/database"
	m "go-fiber-test/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetDogsEndpoint(c *fiber.Ctx) error {
	var dogs []m.Dogs
	if err := database.DBConn.Find(&dogs).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch dogs"})
	}
	return c.Status(fiber.StatusOK).JSON(dogs)
}

func GetDogEndpoint(c *fiber.Ctx) error {
	search := strings.TrimSpace(c.Query("search"))
	if search == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "search is required"})
	}

	var dog m.Dogs
	if err := database.DBConn.Where("dog_id = ?", search).First(&dog).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Dog not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch dog"})
	}

	return c.Status(fiber.StatusOK).JSON(dog)
}

func AddDogEndpoint(c *fiber.Ctx) error {
	var dog m.Dogs
	if err := c.BodyParser(&dog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := database.DBConn.Create(&dog).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create dog"})
	}

	return c.Status(fiber.StatusCreated).JSON(dog)
}

func UpdateDogEndpoint(c *fiber.Ctx) error {
	id := c.Params("id")

	var dog m.Dogs
	if err := c.BodyParser(&dog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	result := database.DBConn.Model(&m.Dogs{}).Where("id = ?", id).Updates(&dog)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update dog"})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Dog not found"})
	}

	return c.Status(fiber.StatusOK).JSON(dog)
}

func RemoveDogEndpoint(c *fiber.Ctx) error {
	id := c.Params("id")

	result := database.DBConn.Delete(&m.Dogs{}, id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete dog"})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Dog not found"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// 7.2 สร้างข้อมูลในตาราง dog มากกว่า 10 ตัว (api add dog) GetdogJson 
func GetDogsJsonV2Endpoint(c *fiber.Ctx) error {
	var dogs []m.Dogs

	if err := database.DBConn.Find(&dogs).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "Could not fetch dogs")
	}

	return c.Status(fiber.StatusOK).JSON(buildDogsResult(dogs))
}

// *7.0.2 สร้าง api GET ใน group dogs โชว์ข้อมูลที่ถูกลบไปแล้ว ตารางdogs
func GetDeletedDogsEndpoint(c *fiber.Ctx) error {
	var dogs []m.Dogs
	if err := database.DBConn.Scopes(m.DeletedDogsScope).Find(&dogs).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "Could not fetch deleted dogs")
	}

	return c.Status(fiber.StatusOK).JSON(dogs)
}

func GetDogsRangeEndpoint(c *fiber.Ctx) error {
	var dogs []m.Dogs

	if err := database.DBConn.Scopes(m.DogIDRangeScope(50, 100)).Find(&dogs).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "Could not fetch dogs")
	}

	return c.Status(fiber.StatusOK).JSON(dogs)
}

// 7.2 สร้างข้อมูลในตาราง dog มากกว่า 10 ตัว (api add dog) GetdogJson 
func GetDogsJson(c *fiber.Ctx) error {
	var dogs []m.Dogs
	if err := database.DBConn.Find(&dogs).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "Could not fetch dogs")
	}

	return c.Status(fiber.StatusOK).JSON(buildDogsResult(dogs))
}

func SeedDogsEndpoint(c *fiber.Ctx) error {
	if err := database.DBConn.Unscoped().Where("1 = 1").Delete(&m.Dogs{}).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, "Could not clear dogs")
	}

	rand.Seed(time.Now().UnixNano())
	colors := []string{"red", "green", "pink", "no color"}
	testDogs := make([]m.Dogs, 0, 15)

	for i := 1; i <= 15; i++ {
		chosenColor := colors[rand.Intn(len(colors))]
		testDogs = append(testDogs, m.Dogs{
			Name:  fmt.Sprintf("Dog-%d", i),
			DogID: generateDogID(chosenColor),
			Color: chosenColor,
		})
	}

	if err := database.DBConn.Create(&testDogs).Error; err != nil {
		return respondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Seed data created successfully",
		"count":   len(testDogs),
		"data":    testDogs,
	})
}
