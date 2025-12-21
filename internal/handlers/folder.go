package handlers

import (
	"strings"

	"github.com/Cypher012/OrganizeNoteAPi/internal/config"
	"github.com/Cypher012/OrganizeNoteAPi/internal/models"
	"github.com/Cypher012/OrganizeNoteAPi/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetFoldersHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := config.GetUserID(c)
		if err != nil {
			return jsonError(c, fiber.StatusUnauthorized, err)
		}

		folders, err := services.GetFolders(db, userId)
		if err != nil {
			return jsonError(c, fiber.StatusInternalServerError, err)
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

		folderSlug := c.Params("slug")
		if folderSlug == "" {
			return jsonError(c, fiber.StatusBadRequest, fiber.ErrBadRequest)
		}

		folder, err := services.GetFolder(db, userId, folderSlug)
		if err != nil {
			switch err {
			case services.ErrFolderNotFound:
				return jsonError(c, fiber.StatusNotFound, err)
			default:
				return jsonError(c, fiber.StatusInternalServerError, err)
			}
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

		var body models.CreateFolderRequest
		if err := c.BodyParser(&body); err != nil {
			return jsonError(c, fiber.StatusBadRequest, err)
		}

		body.Name = strings.TrimSpace(body.Name)
		if err := config.Validate.Struct(body); err != nil {
			return jsonError(c, fiber.StatusBadRequest, err)
		}

		err = services.CreateFolder(db, &body, userId)
		if err != nil {
			switch err {
			case services.ErrFolderExists:
				return jsonError(c, fiber.StatusConflict, err)
			default:
				return jsonError(c, fiber.StatusInternalServerError, err)
			}
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Folder created",
		})
	}
}

func UpdateFolderHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := config.GetUserID(c)
		if err != nil {
			return jsonError(c, fiber.StatusUnauthorized, err)
		}

		body := new(models.UpdateFolderRequest)

		if err := c.BodyParser(body); err != nil {
			return jsonError(c, fiber.StatusBadRequest, err)
		}

		body.Name = strings.TrimSpace(body.Name)
		if err := config.Validate.Struct(body); err != nil {
			return jsonError(c, fiber.StatusBadRequest, err)
		}

		folderSlug := c.Params("slug")
		if folderSlug == "" {
			return jsonError(c, fiber.StatusBadRequest, fiber.ErrBadRequest)
		}

		if err := services.UpdateFolder(db, body, userId, folderSlug); err != nil {
			switch err {
			case services.ErrFolderNotFound:
				return jsonError(c, fiber.StatusNotFound, err)
			default:
				return jsonError(c, fiber.StatusInternalServerError, err)
			}
		}

		return c.JSON(fiber.Map{
			"message": "Folder updated",
		})
	}
}

func DeleteFolderHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := config.GetUserID(c)
		if err != nil {
			return jsonError(c, fiber.StatusUnauthorized, err)
		}

		folderSlug := c.Params("slug")
		if folderSlug == "" {
			return jsonError(c, fiber.StatusBadRequest, fiber.ErrBadRequest)
		}

		err = services.DeleteFolder(db, userId, folderSlug)
		if err != nil {
			switch err {
			case services.ErrFolderNotFound:
				return jsonError(c, fiber.StatusNotFound, err)
			default:
				return jsonError(c, fiber.StatusInternalServerError, err)
			}
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Folder deleted",
		})
	}
}
