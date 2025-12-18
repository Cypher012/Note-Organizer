package handlers

import (
	"errors"
	"strings"

	"github.com/Cypher012/OrganizeNoteAPi/internal/config"
	"github.com/Cypher012/OrganizeNoteAPi/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetFoldersHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := config.GetUserID(c)
		if err != nil {
			return jsonError(c, fiber.StatusUnauthorized, err)
		}

		var folders []models.FolderResponse
		if err := db.
			Model(&models.Folder{}).
			Select("id", "name").
			Where("user_id = ?", userId).
			Find(&folders).Error; err != nil {
			return jsonError(c, fiber.StatusBadRequest, err)
		}

		return c.JSON(fiber.Map{
			"folders": folders,
			"count":   len(folders),
		})
	}
}

func GetFolderByIDHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := config.GetUserID(c)
		if err != nil {
			return jsonError(c, fiber.StatusUnauthorized, err)
		}

		folderId := c.Params("id")
		if folderId == "" {
			return jsonError(c, fiber.StatusBadRequest, fiber.ErrBadRequest)
		}

		var folder models.FolderResponse

		if err := db.
			Model(&models.Folder{}).
			Select("id", "name").
			Where("id = ? AND user_id = ?", folderId, userId).
			First(&folder).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return jsonError(c, fiber.StatusNotFound, fiber.ErrNotFound)
			}

			return jsonError(c, fiber.StatusInternalServerError, err)
		}

		return c.JSON(folder)
	}
}

func CreateFolderHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := config.GetUserID(c)
		if err != nil {
			return jsonError(c, fiber.StatusUnauthorized, err)
		}

		body := new(models.CreateFolderRequest)

		if err := c.BodyParser(body); err != nil {
			return jsonError(c, fiber.StatusBadRequest, err)
		}

		body.Name = strings.TrimSpace(body.Name)
		if err := config.Validate.Struct(body); err != nil {
			return jsonError(c, fiber.StatusBadRequest, err)
		}

		dbFolder := new(models.Folder)
		if err := db.First(dbFolder, "name = ? AND user_id = ?", body.Name, userId).Error; err == nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return jsonError(c, fiber.StatusNotFound, fiber.ErrNotFound)
			}
			return jsonError(c, fiber.StatusInternalServerError, fiber.ErrInternalServerError)
		}

		newFolder := &models.Folder{
			ID:     uuid.NewString(),
			UserID: userId,
			Name:   body.Name,
		}

		if err := db.Create(newFolder).Error; err != nil {
			return jsonError(c, fiber.StatusInternalServerError, err)
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Folder created",
		})
	}
}
