package main

import (
	"log"

	"github.com/Ashu-300/golang-fiber-mvc-starter/database"
	"github.com/Ashu-300/golang-fiber-mvc-starter/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Connect to the database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	app := fiber.New()
	app.Use(logger.New())

	// Setup routes
	routes.Setup(app)

	log.Fatal(app.Listen(":3000"))
}
