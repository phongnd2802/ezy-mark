package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phongnd2802/ezy-mark/internal/models"
	"github.com/phongnd2802/ezy-mark/internal/pkg/context"
	"github.com/phongnd2802/ezy-mark/internal/response"
	"github.com/phongnd2802/ezy-mark/internal/services"
)

type cUserInfo struct{}

var UserInfo = new(cUserInfo)



// GetInfo godoc
// @Summary      Retrieve User Information
// @Description  Fetch the authenticated user's profile information.
// @Tags         UserInfo Management
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Success      200  {object}  response.Response{data=models.UserProfile}
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /user/get-info [get]
func (c *cUserInfo) GetUserProfile(ctx *fiber.Ctx) error {
	sUUID, _ := context.GetSubjectUUID(ctx)
	userId, _ := context.GetUserIdFromUUID(ctx)

	params := models.GetProfileParams{
		UserId: userId,
		SubToken: sUUID,
	}

	code, data, err := services.UserInfo().GetUserProfile(ctx.UserContext(), &params)
	if err != nil {
		return response.ErrorResponse(ctx, code, err)
	}
	return response.SuccessResponse(ctx, code, data)
}