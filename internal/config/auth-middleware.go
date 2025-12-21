package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	models "github.com/Cypher012/OrganizeNoteAPi/internal/models"
	"github.com/Cypher012/OrganizeNoteAPi/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func AuthMiddleWare(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		accessToken := getTokenFromRequest(c)

		if accessToken == "" {
			log.Println("No access token")
			if c.Cookies("rtk") != "" {
				return handleTokenRefresh(c, db)
			}
			return JsonError(c, fiber.StatusUnauthorized, fiber.ErrUnauthorized)
		}

		token, err := jwt.Parse(accessToken, GetJWTSecretKey)

		// Valid access token
		if err == nil && token.Valid {
			userId, err := extractUserId(token)
			if err != nil {
				return JsonError(c, fiber.StatusUnauthorized, fiber.ErrUnauthorized)
			}
			c.Locals("userId", userId)
			return c.Next()
		}

		// Access token exists but is not expired → reject
		if !errors.Is(err, jwt.ErrTokenExpired) {
			return JsonError(c, fiber.StatusUnauthorized, fiber.ErrUnauthorized)
		}

		// Access token expired → refresh
		return handleTokenRefresh(c, db)
	}
}

func getTokenFromRequest(c *fiber.Ctx) string {
	tokenString := c.Cookies("atk")
	if tokenString != "" {
		return tokenString
	}

	authHeader := c.Get("Authorization")
	if after, ok := strings.CutPrefix(authHeader, "Bearer "); ok {
		return after
	}

	return ""
}

func GetJWTSecretKey(t *jwt.Token) (any, error) {
	if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
		return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("critical: JWT_SECRET environment variable is not set")
	}

	return []byte(jwtSecret), nil
}

func extractUserId(token *jwt.Token) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims type")
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("missing or invalid user_id claim")
	}

	return userId, nil
}

func handleTokenRefresh(c *fiber.Ctx, db *gorm.DB) error {
	log.Println("Triggered....")
	refreshToken := c.Cookies("rtk")
	if refreshToken == "" {
		return JsonError(c, fiber.StatusUnauthorized, fiber.ErrUnauthorized)
	}

	rt, err := jwt.Parse(refreshToken, GetJWTSecretKey)
	if err != nil || !rt.Valid {
		return JsonError(c, fiber.StatusUnauthorized, fiber.ErrUnauthorized)
	}

	userId, err := extractUserId(rt)
	if err != nil {
		return JsonError(c, fiber.StatusUnauthorized, fiber.ErrUnauthorized)
	}

	// Verify user exists
	var dbUser models.User
	if err := db.Select("id").First(&dbUser, "id = ?", userId).Error; err != nil {
		return JsonError(c, fiber.StatusUnauthorized, fiber.ErrUnauthorized)
	}

	// Generate new token pair
	newAtk, newRtk, err := services.GenerateJWT(&dbUser)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	SetAuthCookies(c, &newAtk, &newRtk)
	c.Locals("userId", userId)
	return c.Next()
}
