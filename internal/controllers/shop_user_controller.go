package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phongnd2802/ezy-mark/internal/models"
	"github.com/phongnd2802/ezy-mark/internal/pkg/context"
	"github.com/phongnd2802/ezy-mark/internal/response"
	"github.com/phongnd2802/ezy-mark/internal/services"
	"github.com/rs/zerolog/log"
)

type cShopUser struct{}

// Management instance of ShopUser
var ShopUser = new(cShopUser)

// RegisterShop godoc
// @Summary Register a new shop
// @Description Allows users to register a new shop with the required information, including shop name, shop logo, business license, and address.
// @Tags Shop Management
// @Accept multipart/form-data
// @Produce json
// @Param        Authorization header string true "Bearer token"
// @Param        shop_name formData string true "Shop name"
// @Param		 shop_description formData string true "Shop description"
// @Param        shop_logo formData file true "Shop logo"
// @Param 	  	 shop_phone formData string true "Phone number"
// @Param 	 	 shop_email formData string true "Email"
// @Param        shop_address formData string true "Address"
// @Param        business_license formData file true "Business license"
// @Param        tax_id formData string true "Tax ID"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Invalid parameters"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /shop/register [post]
func (c *cShopUser) RegisterShop(ctx *fiber.Ctx) error {
	var params models.RegisterShopReq
	if err := ctx.BodyParser(&params); err != nil {
		return response.ErrorResponse(ctx, response.ErrCodeInvalidParams, err)
	}

	// Get ShopLogo from form-data
	shopLogo, err := ctx.FormFile("shop_logo")
	if err != nil {
		return response.ErrorResponse(ctx, response.ErrCodeInvalidParams, err)
	}
	businessLicense, err := ctx.FormFile("business_license")
	if err != nil {
		return response.ErrorResponse(ctx, response.ErrCodeInvalidParams, err)
	}

	params.ShopLogo = shopLogo
	params.BusinessLicense = businessLicense

	userId, _ := context.GetUserIdFromUUID(ctx)
	log.Info().Msgf("User ID: %d", userId)
	params.OwnerId = userId

	code, err := services.ShopUser().RegisterShop(ctx.UserContext(), &params)
	if err != nil {
		return response.ErrorResponse(ctx, code, err)
	}
	return response.SuccessResponse(ctx, response.ErrCodeSuccess, nil)
}
