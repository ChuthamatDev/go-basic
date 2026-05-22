package routes

import (
	c "go-fiber-test/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

func createBasicAuth(users map[string]string) fiber.Handler {
	return basicauth.New(basicauth.Config{
		Users: users,
		Unauthorized: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized access",
			})
		},
	})
}

func InetRoutes(app *fiber.App) {
	api := app.Group("/api")

	authMiddleware := createBasicAuth(map[string]string{
		"gofiber": "21022566",
	})

	// workshop final
	basicAuth := createBasicAuth(map[string]string{
		"testgo": "23012023",
	})


	// v1 Routes
	v1 := api.Group("/v1")
	v1.Post("/register", c.RegisterEndpoint)

	// แยกกลุ่มที่ใช้ gofiber auth ออกมาให้ชัดเจน (ไม่ใช้ "/" เพื่อไม่ให้ทับกับกลุ่มอื่น)
	// หรือใช้ v1 โดยตรงแล้วใส่ middleware เป็นรายตัวในกลุ่มที่ต้องการ
	
	v1.Get("/", authMiddleware, c.HelloTest)
	v1.Post("/", authMiddleware, c.BodyPersonTest)
	v1.Get("/user-params/:name", authMiddleware, c.ParamsTest) // เปลี่ยนชื่อเพื่อไม่ให้สับสนกับ /user
	v1.Get("/fact/:number", authMiddleware, c.FactorialEndpoint)
	v1.Post("/inet", authMiddleware, c.QueryTest)

	dog := v1.Group("/dog", authMiddleware)
	dog.Get("", c.GetDogsEndpoint)
	dog.Get("/filter", c.GetDogEndpoint)
	dog.Get("/json", c.GetDogsJson)

	// 7.2  สร้างข้อมูลในตาราง dog มากกว่า 10 ตัว (api add dog) GetdogJson
	dog.Get("/json-v2", c.GetDogsJsonV2Endpoint)
	
	dog.Post("/seed", c.SeedDogsEndpoint)
	dog.Get("/deleted", c.GetDeletedDogsEndpoint)
	dog.Get("/range", c.GetDogsRangeEndpoint)
	dog.Post("/", c.AddDogEndpoint)
	dog.Put("/:id", c.UpdateDogEndpoint)
	dog.Delete("/:id", c.RemoveDogEndpoint)

	// Workshop filnal
	userGroup := v1.Group("/user")
	userGroup.Get("", c.GetProfileUserEndpoint)
	userGroup.Get("/generations", c.GetUserGenerationsEndpoint)
	userGroup.Get("/search", c.SearchProfileUsersEndpoint)
	userGroup.Post("/seed", basicAuth, c.SeedProfileUsersEndpoint)
	userGroup.Post("/",basicAuth, c.AddProfileUserEndpoint)
	userGroup.Put("/:id",basicAuth, c.UpdateProfileUserEndpoint)
	userGroup.Delete("/:id",basicAuth, c.RemoveProfileUserEndpoint)

	company := v1.Group("/company", authMiddleware)
	company.Get("/", c.GetCompaniesEndpoint)
	company.Get("/:id", c.GetCompanyByIdEndpoint)
	company.Post("/", c.AddCompanyEndpoint)
	company.Put("/:id", c.UpdateCompanyEndpoint)
	company.Delete("/:id", c.RemoveCompanyEndpoint)

	// v2 Routes
	v2 := api.Group("/v2", authMiddleware)
	v2.Get("/", c.HelloTestV2)

	// v3 Routes
	v3 := api.Group("/v3", authMiddleware)
	v3.Get("/:name", c.AsciiEndpoint)
}