package handlers

import (
	"strings"

	"github.com/Cypher012/OrganizeNoteAPi/internal/config"
	models "github.com/Cypher012/OrganizeNoteAPi/internal/models"
	"github.com/Cypher012/OrganizeNoteAPi/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetNotesHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := config.GetUserID(c)
		if err != nil {
			return jsonError(c, fiber.StatusUnauthorized, err)
		}

		notes, err := services.GetNotes(db, userId)

		if err != nil {
			return jsonError(c, fiber.StatusInternalServerError, fiber.ErrInternalServerError)
		}

		data := fiber.Map{
			"count": len(notes),
			"notes": notes,
		}

		return c.Status(fiber.StatusOK).JSON(data)
	}
}

func GetNotesInFolderHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := config.GetUserID(c)
		if err != nil {
			return jsonError(c, fiber.StatusUnauthorized, err)
		}

		folderSlug := c.Params("folderSlug")

		if folderSlug == "" {
			return jsonError(c, fiber.StatusBadRequest, fiber.ErrBadRequest)

		}

		notes, err := services.GetNotesInFolder(db, userId, folderSlug)

		if err != nil {
			return jsonError(c, fiber.StatusInternalServerError, fiber.ErrInternalServerError)
		}

		data := fiber.Map{
			"count": len(notes),
			"notes": notes,
		}

		return c.Status(fiber.StatusOK).JSON(data)
	}
}

func GetNoteHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := config.GetUserID(c)
		if err != nil {
			return jsonError(c, fiber.StatusUnauthorized, err)
		}

		folderSlug := c.Params("folderSlug")
		noteSlug := c.Params("noteSlug")

		if folderSlug == "" || noteSlug == "" {
			return jsonError(c, fiber.StatusBadRequest, fiber.ErrBadRequest)
		}

		note, err := services.GetNote(db, userId, folderSlug, noteSlug)
		if err != nil {
			switch err {
			case services.ErrNoteNotFound:
				return jsonError(c, fiber.StatusNotFound, err)
			default:
				return jsonError(c, fiber.StatusInternalServerError, fiber.ErrInternalServerError)
			}
		}

		return c.JSON(note)
	}
}

func CreateNoteHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := config.GetUserID(c)
		if err != nil {
			return jsonError(c, fiber.StatusUnauthorized, err)
		}

		folderSlug := c.Params("folderSlug")
		if folderSlug == "" {
			return jsonError(c, fiber.StatusBadRequest, fiber.ErrBadRequest)
		}

		var body models.CreateNoteRequest
		if err := c.BodyParser(&body); err != nil {
			return jsonError(c, fiber.StatusBadRequest, err)
		}

		body.Name = strings.TrimSpace(body.Name)
		if err := config.Validate.Struct(body); err != nil {
			return jsonError(c, fiber.StatusBadRequest, err)
		}

		err = services.CreateNote(db, &body, userId, folderSlug)
		if err != nil {
			switch err {
			case services.ErrNoteExists:
				return jsonError(c, fiber.StatusConflict, err)
			default:
				return jsonError(c, fiber.StatusInternalServerError, fiber.ErrInternalServerError)
			}
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Note created",
		})
	}
}

func UpdateNoteHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := config.GetUserID(c)
		if err != nil {
			return jsonError(c, fiber.StatusUnauthorized, err)
		}

		folderSlug := c.Params("folderSlug")
		noteSlug := c.Params("noteSlug")

		if folderSlug == "" || noteSlug == "" {
			return jsonError(c, fiber.StatusBadRequest, fiber.ErrBadRequest)
		}

		var body models.UpdateNoteRequest
		if err := c.BodyParser(&body); err != nil {
			return jsonError(c, fiber.StatusBadRequest, err)
		}

		body.Name = strings.TrimSpace(body.Name)
		if err := config.Validate.Struct(body); err != nil {
			return jsonError(c, fiber.StatusBadRequest, err)
		}

		err = services.UpdateNote(db, &body, userId, folderSlug, noteSlug)
		if err != nil {
			switch err {
			case services.ErrNoteNotFound:
				return jsonError(c, fiber.StatusNotFound, err)
			default:
				return jsonError(c, fiber.StatusInternalServerError, fiber.ErrInternalServerError)
			}
		}

		return c.JSON(fiber.Map{
			"message": "Note updated",
		})
	}
}

func DeleteNoteHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, err := config.GetUserID(c)
		if err != nil {
			return jsonError(c, fiber.StatusUnauthorized, err)
		}

		folderSlug := c.Params("folderSlug")
		noteSlug := c.Params("noteSlug")

		if folderSlug == "" || noteSlug == "" {
			return jsonError(c, fiber.StatusBadRequest, fiber.ErrBadRequest)
		}

		err = services.DeleteNote(db, userId, folderSlug, noteSlug)
		if err != nil {
			switch err {
			case services.ErrNoteNotFound:
				return jsonError(c, fiber.StatusNotFound, err)
			default:
				return jsonError(c, fiber.StatusInternalServerError, fiber.ErrInternalServerError)
			}
		}

		return c.JSON(fiber.Map{
			"message": "Note deleted",
		})
	}
}
