package utils

import "github.com/gofiber/fiber/v2"

func GetUserID(c *fiber.Ctx) (string, error) {
	userId, ok := c.Locals("userId").(string)
	if !ok || userId == "" {
		return "", fiber.ErrUnauthorized
	}
	return userId, nil
}
