package handlers

import (
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)



func extractNameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 1 {
		return parts[0]
	}
	return ""
}

func getUserKeyOtp(key string) string {
	return fmt.Sprintf("user:%s:otp", key)
}
func getUserKeySession(token string) string {
	return fmt.Sprintf("user:%s:session", token)
}

func render(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/html; charset=utf-8")
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
