package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"flutter_go_crud/backend/internal/handlers"
)

func Register(app *fiber.App, db *gorm.DB) {
	items := handlers.NewItemHandler(db)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "online",
			"message": "Flutter Go CRUD API is running",
			"endpoints": fiber.Map{
				"health":  "/healthz",
				"metrics": "/metrics",
				"api":     "/api/items",
			},
		})
	})

	api := app.Group("/api")
	api.Get("/items", items.List)
	api.Get("/items/:id", items.Get)
	api.Post("/items", items.Create)
	api.Put("/items/:id", items.Update)
	api.Delete("/items/:id", items.Delete)
	api.Delete("/items", items.BulkDelete)
	api.Post("/items/:id/restore", items.Restore)
}
