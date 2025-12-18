package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterNoteRoutes(router fiber.Router, db *gorm.DB) {
	// notes := router.Group("/notes")

	// notes.Post("/", createNote)
	// notes.Get("/:id", getNote)
}
