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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch dogs"})
	}

	return c.Status(fiber.StatusOK).JSON(buildDogsResult(dogs))
}

// *7.0.2 สร้าง api GET ใน group dogs โชว์ข้อมูลที่ถูกลบไปแล้ว ตารางdogs
func GetDeletedDogsEndpoint(c *fiber.Ctx) error {
	var dogs []m.Dogs
	if err := database.DBConn.Scopes(m.DeletedDogsScope).Find(&dogs).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch deleted dogs"})
	}

	return c.Status(fiber.StatusOK).JSON(dogs)
}

func GetDogsRangeEndpoint(c *fiber.Ctx) error {

	var dogs []m.Dogs
	
	if err := database.DBConn.Scopes(m.DogIDRangeScope(50, 100)).Find(&dogs).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch dogs"})
	}

	return c.Status(fiber.StatusOK).JSON(dogs)
}

// 7.2 สร้างข้อมูลในตาราง dog มากกว่า 10 ตัว (api add dog) GetdogJson 
func GetDogsJson(c *fiber.Ctx) error {
	// 1. ดึงข้อมูลสุนัขทั้งหมดจาก Database
	var dogs []m.Dogs
	if err := database.DBConn.Find(&dogs).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch dogs"})
	}

	// 2. เรียกใช้ฟังก์ชันแยก (Helper Function) เพื่อประกอบร่างข้อมูลและนับยอด
	result := buildDogsResult(dogs)

	// 3. ส่งผลลัพธ์กลับเป็น JSON
	return c.Status(fiber.StatusOK).JSON(result)
}

func SeedDogsEndpoint(c *fiber.Ctx) error {
	if err := database.DBConn.Unscoped().Where("1 = 1").Delete(&m.Dogs{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not clear dogs"})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Seed data created successfully",
		"count":   len(testDogs),
		"data":    testDogs,
	})
}

// ฟังก์ชันแยกสำหรับทำ Business Logic: แยกสีและรวมยอด
func buildDogsResult(dogs []m.Dogs) m.ResultDataV3 {
	// ประกาศตัวแปรรอรับผลลัพธ์
	result := m.ResultDataV3{
		Data:       make([]m.DogsRes, 0, len(dogs)),
		Name:       "golang-test",
		Count:      len(dogs),
		SumRed:     0,
		SumGreen:   0,
		SumPink:    0,
		SumNoColor: 0,
	}

	for _, dog := range dogs {
		var color string

		switch {
		case dog.DogID >= 10 && dog.DogID <= 50:
			color = "red"
			result.SumRed++
		case dog.DogID >= 100 && dog.DogID <= 150:
			color = "green"
			result.SumGreen++
		case dog.DogID >= 200 && dog.DogID <= 250:
			color = "pink"
			result.SumPink++
		default:
			color = "no color"
			result.SumNoColor++
		}

		result.Data = append(result.Data, m.DogsRes{
			Name:  dog.Name,
			DogID: dog.DogID,
			Type:  color,
			Color: color, 
		})
	}

	return result
}

func resolveDogColor(dogID int) string {
	switch {
	case dogID >= 10 && dogID <= 50:
		return "red"
	case dogID >= 100 && dogID <= 150:
		return "green"
	case dogID >= 200 && dogID <= 250:
		return "pink"
	default:
		return "no color"
	}
}

func generateDogID(color string) int {
	switch color {
	case "red":
		return rand.Intn(41) + 10
	case "green":
		return rand.Intn(51) + 100
	case "pink":
		return rand.Intn(51) + 200
	default:
		return rand.Intn(9) + 1
	}
}
