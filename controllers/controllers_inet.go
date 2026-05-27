package controllers

import (
	"strconv"
	"strings"

	m "go-fiber-test/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func HelloTest(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func HelloTestV2(c * fiber.Ctx) error {
	return c.SendString("Hello, World! V2")
}

func BodyPersonTest(c *fiber.Ctx) error {
	p := new(m.Person)

	if err := c.BodyParser(p); err != nil {
		return err
	}
	return c.JSON(p)
}

func ParamsTest(c *fiber.Ctx) error {

	str := "hello ==> " + c.Params("name")
	return c.JSON(str)
}

func QueryTest(c *fiber.Ctx) error {

	a := c.Query("search")
	str := "my search is  " + a
	return c.JSON(str)
}

func ValidTest(c *fiber.Ctx) error {
	user := new(m.User)
	if err := c.BodyParser(&user); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Cannot parse JSON")
	}

	if err := validate.Struct(user); err != nil {
		errors := formatValidationErrors(err.(validator.ValidationErrors))
		return respondValidationErrors(c, errors)
	}

	return c.JSON(user)
}	

// 5.1 สร้าง Factorial 
func FactorialEndpoint(c *fiber.Ctx) error {

    param := c.Params("number")

    if param == "" {
        return respondError(c, fiber.StatusBadRequest, "กรุณาระบุตัวเลข")
    }

	// ดักจับขอบเขตข้อมูล
    num, err := strconv.Atoi(param)

    if err != nil || num < 0 {
        return respondError(c, fiber.StatusBadRequest, "กรุณากรอกตัวเลขจำนวนเต็มบวกที่ถูกต้อง")
    }

	//Integer Overflow ค่า 21 แฟกทอเรียลจะมีขนาดบิตที่ใหญ่เกินขอบเขตของประเภทข้อมูล int64
    if num > 20 { 
        return respondError(c, fiber.StatusBadRequest, "ตัวเลขมีขนาดใหญ่เกินไป (รองรับสูงสุดคือ 20)")
    }

    fact := 1
	steps := "1"
    for i := 2; i <= num; i++ {
        fact *= i
		steps = steps + "x" + strconv.Itoa(i)
    }

    return c.JSON(fiber.Map{
        "number":    num,
		"calculation": steps,
        "factorial": fact,
    })
}

// 5.2 สร้าง api v3 ด้วยการ Query Params
func AsciiEndpoint(c *fiber.Ctx) error {
	
	taxID := c.Query("tax_id")

	if taxID == "" {
		return respondError(c, fiber.StatusBadRequest, "จำเป็นต้องระบุ Query parameter 'tax_id'")
	}

	name := c.Params("name")

	asciiStrings := make([]string, len(taxID))
	for i, ch := range taxID {
		asciiStrings[i] = strconv.Itoa(int(ch))
	}

	finalTaxID := strings.Join(asciiStrings, " ")

	return c.JSON(fiber.Map{
		"name":   name,
		"tax_id": finalTaxID, 
	})
}

// 6. api methos POST v1 Resigter 
func RegisterEndpoint(c *fiber.Ctx) error {

	var req m.Register

	if err := c.BodyParser(&req); err != nil {

		return respondError(c, fiber.StatusBadRequest, "รูปแบบข้อมูล JSON ไม่ถูกต้อง")
	}

	if err := validate.Struct(&req); err != nil {

		errors := formatValidationErrors(err.(validator.ValidationErrors))
		
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "การตรวจสอบข้อมูลล้มเหลว",
			"errors":  errors,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "ลงทะเบียนสมาชิกสำเร็จ",
		"data":    req,
	})
}