package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/phongnd2802/daily-social/views/pages"
)

func HandleHome(c echo.Context) error {
	return render(c, pages.Index())
}
