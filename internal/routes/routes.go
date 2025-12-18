package routes

import (
	"github.com/Cypher012/OrganizeNoteAPi/internal/config"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api")

	// Public routes
	RegisterAuthRoutes(api, db)

	// Protected routes
	protected := api.Group("/", config.AuthMiddleWare(db))
	RegisterUserRoutes(protected, db)
	RegisterFolderRoutes(protected, db)
	RegisterNoteRoutes(protected, db)
}
