package services

import (
	"context"

	"github.com/phongnd2802/ezy-mark/internal/models"
)

type (
	// IUserInfo contains methods for users to update their personal information
	IUserInfo interface {
		GetUserProfile(ctx context.Context, params *models.GetProfileParams) (code int, res *models.UserProfileRes, err error)
		UpdateUserProfile(ctx context.Context, params *models.UpdateUserProfileReq) (code int, res *models.UserProfileRes, err error)
		ChangePassword(ctx context.Context, params *models.ChangePassword) (code int, err error)
	}

	// IUserAdmin contains methods for administrators to manage user accounts
	IUserAdmin interface {
		BlockUser(ctx context.Context) error
		DeleteUser(ctx context.Context) error
	}
)

var (
	localUserInfo  IUserInfo
	localUserAdmin IUserAdmin
)

func UserInfo() IUserInfo {
	if localUserInfo == nil {
		panic("IUserInfo interface not implemented yet")
	}
	return localUserInfo
}

func InitUserInfo(i IUserInfo) {
	localUserInfo = i
}

func UserAdmin() IUserAdmin {
	if localUserAdmin == nil {
		panic("IUserAdmin interface not implemented yet")
	}

	return localUserAdmin
}

func InitUserAdmin(i IUserAdmin) {
	localUserAdmin = i
}
