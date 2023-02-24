package routes

import (
	"github.com/gofiber/fiber/v2"
	"bola-oi/controllers"
)

func FieldRoute(app *fiber.App) {
	app.Post("/api/field", controllers.CreateField)
	app.Put("/api/field/:fieldId", controllers.EditField)
	app.Get("/api/field/:fieldId", controllers.GetFieldById)
	app.Get("/api/field", controllers.GetAllFields)
	app.Delete("/api/field/:fieldId", controllers.DeleteField)
}