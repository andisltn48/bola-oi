package routes

import (
	"github.com/gofiber/fiber/v2"
	"bola-oi/controllers"
)

func AuthRoute(app *fiber.App) {
	app.Post("/api/auth/register", controllers.Register)
	app.Post("/api/auth/login", controllers.Login)
}