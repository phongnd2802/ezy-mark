package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phongnd2802/ezy-mark/internal/pkg/token"
	"github.com/phongnd2802/ezy-mark/internal/response"
)

func AuthenticationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check headers authorization
		jwtToken, exist := token.ExtractBearerToken(c)
		if !exist {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Response{
				Code:    401,
				Message: "unauthorized",
			})
		}
		
		

		// validate jwt token by subject
		claims, err := token.VerifyTokenSubject(jwtToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Response{
				Code:    401,
				Message: "invalid token",
			})
		}

		// Store user ID in request context
		c.Locals("subjectUUID", claims.Subject)
		return c.Next()
	}
}
