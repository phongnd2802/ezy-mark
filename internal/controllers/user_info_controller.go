package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/phongnd2802/ezy-mark/internal/models"
	"github.com/phongnd2802/ezy-mark/internal/pkg/context"
	"github.com/phongnd2802/ezy-mark/internal/response"
	"github.com/phongnd2802/ezy-mark/internal/services"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

type cUserInfo struct{}

var UserInfo = new(cUserInfo)



// ChangePassword godoc
// @Summary      Change Password
// @Description  Allows users to change their password by providing the old password and the new password.
// @Tags         UserInfo Management
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        payload body models.ChangePassword true "payload"     
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /user/change-password [patch]
func (c *cUserInfo) ChangePassword(ctx *fiber.Ctx) error {
	params := new(models.ChangePassword)
	if err := ctx.BodyParser(params); err != nil {
		return response.ErrorResponse(ctx, response.ErrCodeInvalidParams, err)
	}

	userId, _ := context.GetUserIdFromUUID(ctx)
	params.UserId = userId

	code, err := services.UserInfo().ChangePassword(ctx.UserContext(), params)
	if err != nil {
		return response.ErrorResponse(ctx, code, err)
	}
	return response.SuccessResponse(ctx, response.ErrCodeSuccess, nil)
}


// GetInfo godoc
// @Summary      Retrieve User Information
// @Description  Fetch the authenticated user's profile information.
// @Tags         UserInfo Management
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Success      200  {object}  response.Response{data=models.UserProfileRes}
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /user/get-info [get]
func (c *cUserInfo) GetUserProfile(ctx *fiber.Ctx) error {
	sUUID, _ := context.GetSubjectUUID(ctx)
	userId, _ := context.GetUserIdFromUUID(ctx)
	roles := ctx.Locals("roles").([]int32)
	log.Info().Msgf("Roles: %v", roles)
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
// @Param user_avatar formData file false "User Avatar File (Only images: jpg, jpeg, png)"
// @Success 200 {object} response.Response{data=models.UserProfileRes} "Profile updated successfully"
// @Failure 400 {object} response.Response "Invalid parameters"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /user/update-info [patch]
func (c *cUserInfo) UpdateUserProfile(ctx *fiber.Ctx) error {
	params := new(models.UpdateUserProfileReq)

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
	if file != nil {
		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/jpg": true,
			"image/png":  true,
		}

		src, err := file.Open()
		if err != nil {
			return response.ErrorResponse(ctx, response.ErrCodeInvalidParams, err)
		}
		defer src.Close()

		buffer := make([]byte, 512)
		_, err = src.Read(buffer)
		if err != nil {
			return response.ErrorResponse(ctx, response.ErrCodeInvalidParams, err)
		}

		fileType := http.DetectContentType(buffer)
		if !allowedTypes[fileType] {
			return response.ErrorResponse(ctx, response.ErrCodeInvalidParams, errors.New("only image files (jpg, jpeg, png) are allowed"))
		}

		params.UserAvatar = file
	}

	// Get UserID
	userId, _ := context.GetUserIdFromUUID(ctx)
	params.UserId = userId

	// Get SubToken
	subToken, _ := context.GetSubjectUUID(ctx)
	params.SubToken = subToken

	// Call Service to Update User Profile
	code, data, err := services.UserInfo().UpdateUserProfile(ctx.UserContext(), params)
	if err != nil {
		return response.ErrorResponse(ctx, code, err)
	}
	return response.SuccessResponse(ctx, response.ErrCodeSuccess, data)
}
