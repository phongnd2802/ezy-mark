package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/phongnd2802/ezy-mark/internal/models"
	"github.com/phongnd2802/ezy-mark/internal/pkg/context"
	"github.com/phongnd2802/ezy-mark/internal/response"
	"github.com/phongnd2802/ezy-mark/internal/services"
	"github.com/valyala/fasthttp"
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
		UserId:   userId,
		SubToken: sUUID,
	}

	code, data, err := services.UserInfo().GetUserProfile(ctx.UserContext(), &params)
	if err != nil {
		return response.ErrorResponse(ctx, code, err)
	}
	return response.SuccessResponse(ctx, code, data)
}

// UpdateInfo godoc
// @Summary Update user profile
// @Description Allows users to update their profile information, including nickname, full name, mobile, gender, birthday, and avatar.
// @Tags UserInfo Management
// @Accept multipart/form-data
// @Produce json
// @Param        Authorization header string true "Bearer token"
// @Param user_nickname formData string true "User Nickname"
// @Param user_fullname formData string false "User Full Name"
// @Param user_mobile formData string false "User Mobile Number"
// @Param user_gender formData string false "User Gender (male, female, other)"
// @Param user_birthday formData string false "User Birthday (YYYY-MM-DD)"
// @Param user_avatar formData file false "User Avatar File"
// @Success 200 {object} response.Response{data=models.UpdateProfileUserRes} "Profile updated successfully"
// @Failure 400 {object} response.Response "Invalid parameters"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /user/update-info [patch]
func (c *cUserInfo) UpdateUserProfile(ctx *fiber.Ctx) error {
	params := new(models.UpdateProfileUserReq)

	// Parse Body
	if err := ctx.BodyParser(params); err != nil {
		return response.ErrorResponse(ctx, response.ErrCodeInvalidParams, err)
	}

	// validate birthday
	if params.UserBirthday != "" {
		_, err := time.Parse("2006-01-02", params.UserBirthday)
		if err != nil {
			return response.ErrorResponse(ctx, response.ErrCodeInvalidParams, err)
		}
	}

	// Get File Avatar
	file, err := ctx.FormFile("user_avatar")
	if err != nil && err != fasthttp.ErrMissingFile {
		return response.ErrorResponse(ctx, response.ErrCodeInvalidParams, err)
	}
	params.UserAvatar = file

	// Get UserID
	userId, _ := context.GetUserIdFromUUID(ctx)
	params.UserId = userId

	// Get SubToken
	subToken, _ := context.GetSubjectUUID(ctx)
	params.SubToken = subToken

	code, data, err := services.UserInfo().UpdateUserProfile(ctx.UserContext(), params)
	if err != nil {
		return response.ErrorResponse(ctx, code, err)
	}
	return response.SuccessResponse(ctx, response.ErrCodeSuccess, data)
}
