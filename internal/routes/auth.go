package routes

import (
	handlers "github.com/Cypher012/OrganizeNoteAPi/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(router fiber.Router, db *gorm.DB) {
	auth := router.Group("/auth")

	auth.Post("/register", handlers.RegisterHandler(db))
	auth.Post("/login", handlers.LoginHandler(db))
}
