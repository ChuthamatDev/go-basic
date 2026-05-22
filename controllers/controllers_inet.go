package controllers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	m "go-fiber-test/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)


var validate = validator.New()

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	validBizTypes = map[string]struct{}{"E-commerce": {}, "Service": {}, "IT": {}, "Consulting": {}, "Other": {}}
	subdomainRegex = regexp.MustCompile(`^[a-z0-9.-]{2,30}$`)
)

func init() {
	validate.RegisterValidation("username_custom", func(fl validator.FieldLevel) bool {
		return usernameRegex.MatchString(fl.Field().String())
	})

	validate.RegisterValidation("business_type_custom", func(fl validator.FieldLevel) bool {
		_, ok := validBizTypes[fl.Field().String()]
		return ok
	})

	validate.RegisterValidation("website_custom", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()

		if !subdomainRegex.MatchString(value) {
			return false
		}

		allowedSuffixes := []string{
			".sogodweb.com",
			".sogodweb.co.th",
			".sogodweb.in.th",
		}
		for _, suffix := range allowedSuffixes {
			if strings.HasSuffix(value, suffix) {
				return true
			}
		}
		return false
	})
}

func formatValidationErrors(errs validator.ValidationErrors) []string {
	var errorMessages []string
	for _, err := range errs {
		var message string
		switch err.Tag() {
		case "required":
			message = fmt.Sprintf("ฟิลด์ %s เป็นข้อมูลที่จำเป็น", err.Field())
		case "email":
			message = "รูปแบบอีเมลไม่ถูกต้อง"
		case "min":
			message = fmt.Sprintf("ฟิลด์ %s ต้องมีความยาวอย่างน้อย %s ตัวอักษร", err.Field(), err.Param())
		case "max":
			message = fmt.Sprintf("ฟิลด์ %s ต้องมีความยาวไม่เกิน %s ตัวอักษร", err.Field(), err.Param())
		case "username_custom":
			message = "Username ใช้อักษรภาษาอังกฤษ (a-z, A-Z), ตัวเลข (0-9) และเครื่องหมาย (_), (-) เท่านั้น"
		case "business_type_custom":
			message = "ประเภทธุรกิจที่เลือกไม่ถูกต้อง"
		case "website_custom":
			message = "Website ต้องเป็น subdomain ที่ลงท้ายด้วย .sogodweb.com, .sogodweb.co.th, หรือ .sogodweb.in.th"
		case "numeric":
			message = fmt.Sprintf("ฟิลด์ %s ต้องเป็นตัวเลขเท่านั้น", err.Field())
		default:
			message = fmt.Sprintf("ข้อมูลฟิลด์ %s ไม่ถูกต้อง", err.Field())
		}
		errorMessages = append(errorMessages, message)
	}
	return errorMessages
}

func HelloTest(c * fiber.Ctx) error {
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := validate.Struct(user); err != nil {
		errors := formatValidationErrors(err.(validator.ValidationErrors))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Validation failed", "errors": errors})
	}

	return c.JSON(user)
}	

// 5.1 สร้าง Factorial 
func FactorialEndpoint(c *fiber.Ctx) error {

    param := c.Params("number")

    if param == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "กรุณาระบุตัวเลข"})
    }

	// ดักจับขอบเขตข้อมูล
    num, err := strconv.Atoi(param)

    if err != nil || num < 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "กรุณากรอกตัวเลขจำนวนเต็มบวกที่ถูกต้อง"})
    }

	//Integer Overflow ค่า 21 แฟกทอเรียลจะมีขนาดบิตที่ใหญ่เกินขอบเขตของประเภทข้อมูล int64
    if num > 20 { 
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ตัวเลขมีขนาดใหญ่เกินไป (รองรับสูงสุดคือ 20)"})
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "จำเป็นต้องระบุ Query parameter 'tax_id'"})
	}

	name := c.Params("name")
	fmt.Printf("Params 'name': %s\n", name)

	asciiStrings := make([]string, len(taxID))
	fmt.Printf("Slice ก่อนเข้าลูป: %q\n", asciiStrings)

	for i, ch := range taxID {
		asciiStrings[i] = strconv.Itoa(int(ch))
		
		fmt.Printf("   รอบที่ i=%d: ดึงอักษร '%c' -> แปลงเป็นข้อความ \"%s\" -> Slice ปัจจุบัน: %q\n", i, ch, asciiStrings[i], asciiStrings)
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

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "รูปแบบข้อมูล JSON ไม่ถูกต้อง",
		})
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