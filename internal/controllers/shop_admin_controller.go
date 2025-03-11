package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phongnd2802/ezy-mark/internal/models"
	"github.com/phongnd2802/ezy-mark/internal/pkg/context"
	"github.com/phongnd2802/ezy-mark/internal/response"
	"github.com/phongnd2802/ezy-mark/internal/services"
)

type cShopAdmin struct{}

var ShopAdmin = new(cShopAdmin)


// ApproveShop godoc
// @Summary      Approve a shop
// @Description  Allows an admin to approve a shop by providing the shop ID.
// @Tags         ShopAdmin Management
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        payload body models.ApproveShopReq true "payload"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /admin/shop/approve [patch]
func (c *cShopAdmin) ApproveShop(ctx *fiber.Ctx) error {
	var params models.ApproveShopReq
	if err := ctx.BodyParser(&params); err != nil {
		return response.ErrorResponse(ctx, response.ErrCodeInvalidParams, err)
	}
	userId, _ := context.GetUserIdFromUUID(ctx)
	params.UserId = userId

	code, err := services.ShopAdmin().ApproveShop(ctx.UserContext(), &params)
	if err != nil {
		return response.ErrorResponse(ctx, code, err)
	}
	return response.SuccessResponse(ctx, code, nil)
}