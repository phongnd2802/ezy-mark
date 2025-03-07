package middlewares

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/phongnd2802/ezy-mark/internal/pkg/token"
	"github.com/rs/zerolog/log"
)

func AuthenticationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check headers authorization
		jwtToken, exist := token.ExtractBearerToken(c)
		if !exist {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		// validate jwt token by subject
		claims, err := token.VerifyTokenSubject(jwtToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		// update claims to context
		log.Info().Str("claims::: UUID::", claims.Subject).Msg("Authenticated user")

		ctx := context.WithValue(c.UserContext(), "subjectUUID", claims.Subject)
		c.SetUserContext(ctx)

		return c.Next()
	}
}