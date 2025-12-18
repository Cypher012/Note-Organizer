package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterUserRoutes(router fiber.Router, db *gorm.DB) {
	// users := router.Group("/users")

	// users.Get("/", getUsers)
	// users.Get("/:id", getUser)
}
