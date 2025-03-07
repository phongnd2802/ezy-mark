package token

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ExtractBearerToken(c *fiber.Ctx) (string, bool) {
	authHeader := c.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer") {
		return strings.TrimPrefix(authHeader, "Bearer "), true
	}

	return "", false
}

