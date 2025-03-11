package routes

import (
	"github.com/phongnd2802/ezy-mark/internal/routes/admin"
	"github.com/phongnd2802/ezy-mark/internal/routes/user"
)

type RouterGroup struct {
	User  user.UserRouterGroup
	Admin admin.AdminRouterGroup
}

var RouterGroupApp = new(RouterGroup)
