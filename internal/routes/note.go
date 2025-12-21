package routes

import (
	"github.com/Cypher012/OrganizeNoteAPi/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterNoteRoutes(router fiber.Router, db *gorm.DB) {
	router.Get("/notes", handlers.GetNotesHandler(db))

	// notes inside folders
	folders := router.Group("/folders/:folderSlug")

	folders.Get("/notes", handlers.GetNotesInFolderHandler(db))
	folders.Post("/notes", handlers.CreateNoteHandler(db))
	folders.Get("/notes/:noteSlug", handlers.GetNoteHandler(db))
	folders.Put("/notes/:noteSlug", handlers.UpdateNoteHandler(db))
	folders.Delete("/notes/:noteSlug", handlers.DeleteNoteHandler(db))

}
