package routes

import "github.com/phongnd2802/ezy-mark/internal/routes/user"

type RouterGroup struct {
	User user.UserRouterGroup
}

var RouterGroupApp = new(RouterGroup)
