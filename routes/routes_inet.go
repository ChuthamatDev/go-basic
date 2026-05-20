package routes

import (
	"go-fiber-test/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

func InetRoutes(app *fiber.App) {
	app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"gofiber": "21022566",
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized access",
			})
		},
	}))

	// version v1
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v2 := api.Group("/v2")

	v2.Get("/",controllers.HelloTestV2 )
	v1.Get("/",controllers.HelloTest )
	v1.Post("/",controllers.BodyPersonTest )
	v1.Get("/user/:name",controllers.ParamsTest )
	v1.Get("/fact/:number",controllers.Factorial )
	v1.Post("/inet",controllers.QueryTest )
	v1.Post("/valid",controllers.ValidTest )

	//CRUD dogs
   dog := v1.Group("/dog")
   dog.Get("", controllers.GetDogs)
   dog.Get("/filter", controllers.GetDog)
   dog.Get("/json", controllers.GetDogsJson)
   dog.Post("/", controllers.AddDog)
   dog.Put("/:id", controllers.UpdateDog)
   dog.Delete("/:id", controllers.RemoveDog)

   //Factorial
	fact := v1.Group("/fact")
	fact.Get("/:number", controllers.Factorial)
}