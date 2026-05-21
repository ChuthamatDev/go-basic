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

	// v1 Routes
	v1 := api.Group("/v1")
	v1.Post("/register", c.RegisterEndpoint)

	v1Auth := v1.Group("/", authMiddleware)
	v1Auth.Get("/", c.HelloTest)
	v1Auth.Post("/", c.BodyPersonTest)
	v1Auth.Get("/user/:name", c.ParamsTest)
	v1Auth.Get("/fact/:number", c.FactorialEndpoint)
	v1Auth.Post("/inet", c.QueryTest)
	// v1Auth.Post("/valid", c.ValidTest)

	dog := v1Auth.Group("/dog")
	dog.Get("", c.GetDogsEndpoint)
	dog.Get("/filter", c.GetDogEndpoint)
	dog.Get("/json", c.GetDogsJson)
	dog.Get("/json-v2", c.GetDogsJsonV2Endpoint)
	dog.Post("/seed", c.SeedDogsEndpoint)
	dog.Get("/deleted", c.GetDeletedDogsEndpoint)
	dog.Get("/range", c.GetDogsRangeEndpoint)
	dog.Post("/", c.AddDogEndpoint)
	dog.Put("/:id", c.UpdateDogEndpoint)
	dog.Delete("/:id", c.RemoveDogEndpoint)

	// Company Routes
	company := v1Auth.Group("/company")
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