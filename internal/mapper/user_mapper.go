package mapper

import (
	"time"

	"github.com/phongnd2802/ezy-mark/internal/db"
	"github.com/phongnd2802/ezy-mark/internal/models"
)

func MapUserProfile(dbProfile db.GetUserProfileRow) *models.UserProfile {
	var userBirthday *time.Time
	if dbProfile.UserBirthday.Valid{
		userBirthday = &dbProfile.UserBirthday.Time
	}

	var userGender *bool
	if dbProfile.UserGender.Valid {
		userGender = &dbProfile.UserGender.Bool
	}

	return &models.UserProfile{
		UserNickname: dbProfile.UserNickname,
		UserFullname: dbProfile.UserFullname.String,
		UserAvatar: dbProfile.UserAvatar.String,
		UserMobile: dbProfile.UserMobile.String,
		UserBirthday: userBirthday,
		UserGender: userGender,
	}
}