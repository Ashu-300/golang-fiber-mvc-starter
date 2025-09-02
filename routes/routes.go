package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/Ashu-300/golang-fiber-mvc-starter/controllers"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Get("/user", controllers.GetUser)
	api.Post("/logout", controllers.Logout)
}

