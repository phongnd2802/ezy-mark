package mapper

import (
	"github.com/phongnd2802/ezy-mark/internal/db"
	"github.com/phongnd2802/ezy-mark/internal/models"
)

func MapUserProfile(dbProfile db.GetUserProfileRow) *models.UserProfileRes {
	var userGender *string
	switch dbProfile.UserGender.GenderEnum {
	case "male", "female", "other":
		userGender = (*string)(&dbProfile.UserGender.GenderEnum)
	default:
		userGender = nil
	}

	return &models.UserProfileRes{
		UserProfile: models.UserProfile{
			UserNickname: dbProfile.UserNickname,
			UserFullname: dbProfile.UserFullname.String,
			UserMobile:   dbProfile.UserMobile.String,
			UserBirthday: func () string {
				if dbProfile.UserBirthday.Valid {
					return dbProfile.UserBirthday.Time.String()
				}
				return ""
			}(),
			UserGender:   userGender,
		},
		UserAvatar: dbProfile.UserAvatar.String,
	}
}
