package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var Validate = validator.New()

func JsonError(c *fiber.Ctx, status int, err error) error {
	return c.Status(status).JSON(fiber.Map{
		"error": err.Error(),
	})
}

func GetUserID(c *fiber.Ctx) (userId string, err error) {
	userId, ok := c.Locals("userId").(string)
	if !ok || userId == "" {
		return "", fiber.ErrUnauthorized
	}
	return userId, nil
}

func SetAuthCookies(c *fiber.Ctx, atk, rtk *string) {
	if atk != nil {
		c.Cookie(&fiber.Cookie{
			Name:     "atk",
			Value:    *atk,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
			Path:     "/",
			MaxAge:   60 * 30, // 30 minutes
		})
	}

	if rtk != nil {
		c.Cookie(&fiber.Cookie{
			Name:     "rtk",
			Value:    *rtk,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
			Path:     "/",
			MaxAge:   3600 * 24 * 7, // 7 days
		})
	}
}

func ClearAuthCookies(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:   "atk",
		MaxAge: -1,
		Path:   "/",
	})
	c.Cookie(&fiber.Cookie{
		Name:   "rtk",
		MaxAge: -1,
		Path:   "/",
	})
}
