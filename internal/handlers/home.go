package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/phongnd2802/daily-social/views/components"
	"github.com/phongnd2802/daily-social/views/pages"
)

func (h *Handler) HandleHome(c echo.Context) error {
	method := c.Request().Method
	if method == echo.GET {
		return render(c, pages.Index())
	}
	
	return render(c, components.Toast(components.ToastProps{
		Message: "Nguyen Duy Phong",
		Type: "error",
	}))
}
