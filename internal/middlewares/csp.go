package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/phongnd2802/daily-social/internal/pkg/utils"
)

type CSPConfig struct {
	ScriptSrc []string // External script domains allowed
}

func WithCSP(config CSPConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			nonce, err := utils.GenerateNonce()
			if err != nil {
				log.Printf("failed to generate nonce: %v", err)
				c.Response().Header().Set("Content-Security-Policy", "script-src 'self'")
				return c.NoContent(http.StatusInternalServerError)
			}

			// Combine all script sources
			scriptSources := append(
				[]string{"'self'", fmt.Sprintf("'nonce-%s'", nonce)},
				config.ScriptSrc...)
			csp := fmt.Sprintf("script-src %s", strings.Join(scriptSources, " "))
			c.Response().Header().Set("Content-Security-Policy", csp)
			return next(c)
		}
	}
}
