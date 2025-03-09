package context

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/phongnd2802/ezy-mark/internal/helpers"
	"github.com/phongnd2802/ezy-mark/internal/pkg/cache"
)

type UserInfoUUID struct {
	UserId    int64  `json:"user_id"`
	UserEmail string `json:"user_email"`
}

func GetSubjectUUID(ctx *fiber.Ctx) (string, error) {
	sUUID, ok := ctx.Locals("subjectUUID").(string)
	if !ok {
		return "", fmt.Errorf("failed to get subject UUID")
	}
	return sUUID, nil
}

func GetUserIdFromUUID(ctx *fiber.Ctx) (int64, error) {
	sUUID, err := GetSubjectUUID(ctx)
	log.Info("subjectUUID: ", sUUID)
	if err != nil {
		return 0, nil
	}

	// Get UserInfo from redis by uuid
	var userInfo UserInfoUUID
	if err := cache.GetCache(ctx.UserContext(), helpers.GetUserKeyProfile(sUUID), &userInfo); err != nil {
		log.Info("error: ", err)
		return 0, err
	}
	log.Info("userInfo: ", userInfo)
	return userInfo.UserId, nil
}


func GetRoles(ctx *fiber.Ctx) ([]string, error) {
	roles, ok := ctx.Locals("roles").([]string)
	if !ok {
		return nil, fmt.Errorf("failed to get roles")
	}
	return roles, nil
}

func GetPermissions(ctx *fiber.Ctx) ([]string, error) {
	permissions, ok := ctx.Locals("permissions").([]string)
	if !ok {
		return nil, fmt.Errorf("failed to get permissions")
	}
	return permissions, nil
}