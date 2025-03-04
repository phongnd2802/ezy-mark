package middlewares

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ContextKey string

const (
	UserAgentKey ContextKey = "UserAgent"
	ClientIPKey  ContextKey = "ClientIP"
)

func RequestMetadataMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userAgent := c.Get("User-Agent")
		clientIP := c.Get("X-Real-IP")
		if clientIP == "" {
			clientIP = c.Get("X-Forwarded-For")
		}

		if clientIP == "" {
			clientIP = c.IP()
		}

		clientIP = strings.Split(clientIP, ",")[0]

		ctx := context.WithValue(c.UserContext(), UserAgentKey, userAgent)
		ctx = context.WithValue(ctx, ClientIPKey, clientIP)

		c.SetUserContext(ctx)

		return c.Next()
	}
}
