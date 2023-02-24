package main

import (
	"github.com/gofiber/fiber/v2"
	"bola-oi/configs"
	"bola-oi/routes"
)

func main() {
	app:= fiber.New()

	configs.ConnectDB()

	routes.FieldRoute(app)
	routes.AuthRoute(app)

	app.Listen(":5000")
}