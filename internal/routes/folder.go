package routes

import (
	"github.com/Cypher012/OrganizeNoteAPi/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterFolderRoutes(router fiber.Router, db *gorm.DB) {
	folders := router.Group("/folder")

	folders.Get("/", handlers.GetFoldersHandler(db))
	folders.Get("/:slug", handlers.GetFolderByIDHandler(db))
	folders.Post("/", handlers.CreateFolderHandler(db))
	folders.Put("/:slug", handlers.UpdateFolderHandler(db))
	folders.Delete("/:slug", handlers.DeleteFolderHandler(db))

	//  notes inside a folder
	folders.Get("/:folderSlug/notes", handlers.GetNotesInFolderHandler(db))
}
