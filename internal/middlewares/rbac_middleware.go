package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phongnd2802/ezy-mark/internal/helpers"
	"github.com/phongnd2802/ezy-mark/internal/pkg/context"
	"github.com/phongnd2802/ezy-mark/internal/response"
)



// RBAC middleware to check if the user has the required roles
func RBAC(roles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoles, err := context.GetRoles(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.Response{
				Code:    500,
				Message: "internal server error",
				Data:    err.Error(),
			})
		}
		hasRole := helpers.HasValidRole(roles, userRoles)
		if !hasRole {
			return c.Status(fiber.StatusForbidden).JSON(response.Response{
				Code:    403,
				Message: "forbidden",
			})
		}

		return c.Next()
	}
}

// PermissionMiddleware to check if the user has the required permissions
func PermissionMiddleware(permissions []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO...
		return c.Next()
	}
}
