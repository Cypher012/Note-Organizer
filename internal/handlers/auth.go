package handlers

import (
	"errors"
	"fmt"

	config "github.com/Cypher012/OrganizeNoteAPi/internal/config"
	models "github.com/Cypher012/OrganizeNoteAPi/internal/models"
	"github.com/Cypher012/OrganizeNoteAPi/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var jsonError = config.JsonError

func RegisterHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body models.RegisterRequest

		if err := c.BodyParser(&body); err != nil {
			return jsonError(c, 400, errors.New("invalid request body"))
		}

		if err := config.Validate.Struct(body); err != nil {
			e := err.(validator.ValidationErrors)[0]
			return config.JsonError(
				c,
				fiber.StatusBadRequest,
				fmt.Errorf("%s is invalid", e.Field()),
			)
		}

		user, accessToken, refreshToken, err := services.RegisterUser(db, &body)

		if err != nil {
			return jsonError(c, fiber.StatusBadRequest, err)
		}

		config.SetAuthCookies(c, &accessToken, &refreshToken)

		return c.JSON(fiber.Map{
			"user": fiber.Map{
				"username": user.Username,
				"email":    user.Email,
				"role":     user.Role,
			},
			"Message": "Registration successful",
		})
	}
}

func LoginHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := new(models.LoginRequest)

		if err := c.BodyParser(body); err != nil {
			return jsonError(c, fiber.StatusBadRequest, fiber.ErrBadRequest)
		}

		user, accessToken, refreshToken, err := services.LoginUser(db, body)

		if err != nil {
			return jsonError(c, fiber.StatusBadRequest, fiber.ErrBadRequest)
		}

		config.SetAuthCookies(c, &accessToken, &refreshToken)

		return c.JSON(fiber.Map{
			"user": fiber.Map{
				"username": user.Username,
				"email":    user.Email,
				"role":     user.Role,
			},
			"message": "Login successful",
		})
	}
}
