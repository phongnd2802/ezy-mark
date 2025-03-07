package context

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
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
	if err != nil {
		return 0, nil
	}

	// Get UserInfo from redis by uuid
	var userInfo UserInfoUUID
	if err := cache.GetCache(ctx.UserContext(), sUUID, &userInfo); err != nil {
		return 0, err
	}
	return userInfo.UserId, nil
}
