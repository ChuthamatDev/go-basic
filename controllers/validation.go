package controllers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
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
