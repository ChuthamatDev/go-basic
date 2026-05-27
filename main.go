package main

import (
	"log"

	"go-fiber-test/database"
	"go-fiber-test/models"
	"go-fiber-test/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	if err := database.Init(); err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}

	if err := database.Migrate(&models.Dogs{}, &models.Company{}, &models.ProfileUser{}); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	routes.InetRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
