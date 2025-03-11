package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/phongnd2802/ezy-mark/internal/global"
	"github.com/phongnd2802/ezy-mark/internal/helpers"
	"github.com/phongnd2802/ezy-mark/internal/pkg/cache"
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

		// Check the blacklist
		isTokenExistsBlackList, err := global.Rdb.Exists(c.UserContext(), helpers.GetKeyBlackList(claims.Subject)).Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
                Code:    500,
                Message: "internal server error",
                Data:    err.Error(),
            })
		}
		if isTokenExistsBlackList > 0 {
            return c.Status(fiber.StatusUnauthorized).JSON(response.Response{
                Code:    401,
                Message: "token revoked",
            })
        }

		// Store user ID in request context
		c.Locals("subjectUUID", claims.Subject)
		log.Info("subjectUUID: ", claims.Subject)
		// Store role in request context
		var roles []string
		err = cache.GetCache(c.UserContext(), helpers.GetUserKeyRole(claims.Subject), &roles)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
				Code:    500,
				Message: "internal server error",
				Data:    err.Error(),
			})
		}

		c.Locals("roles", roles)
		return c.Next()
	}
}
