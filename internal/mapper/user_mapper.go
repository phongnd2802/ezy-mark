package mapper

import (
	"github.com/phongnd2802/ezy-mark/internal/db"
	"github.com/phongnd2802/ezy-mark/internal/models"
)

func MapUserProfile(dbProfile db.GetUserProfileRow) *models.UserProfile {
	var userGender *string
	switch dbProfile.UserGender.GenderEnum {
	case "male", "female", "other":
		userGender = (*string)(&dbProfile.UserGender.GenderEnum)
	default:
		userGender = nil
	}

	userBirthday := dbProfile.UserBirthday.Time.String()
	return &models.UserProfile{
		UserNickname: dbProfile.UserNickname,
		UserFullname: dbProfile.UserFullname.String,
		UserMobile:   dbProfile.UserMobile.String,
		UserBirthday: userBirthday,
		UserGender:   userGender,
	}
}
