package services

import (
	"context"

	"github.com/phongnd2802/ezy-mark/internal/db"
	"github.com/phongnd2802/ezy-mark/internal/models"
)

type (
	IShopUser interface {
		RegisterShop(ctx context.Context, params *models.RegisterShopReq) (code int, err error)
		GetShopDetails(ctx context.Context, shopId int64) (code int, res *models.ShopDetailsRes, err error)
	}

	IShopOwner interface{
		GetShopByOwner(ctx context.Context, ownerId int64) (code int, res []db.Shop, err error)
	}

	IShopAdmin interface {
		ApproveShop(ctx context.Context, params *models.ApproveShopReq) (code int, err error)
		GetAllShops(ctx context.Context) (code int, res []db.Shop, err error)
		GetShopById(ctx context.Context, shopId int64) (code int, res db.Shop, err error)
		BlockShop(ctx context.Context, shopId int64) (code int, err error)
	}
)

var (
	localShopUser IShopUser
	localShopAdmin IShopAdmin
	localShopOwner IShopOwner
)

func ShopUser() IShopUser {
	if localShopUser == nil {
		panic("IShopUser interface not implemented yet")
	}
	return localShopUser
}

func InitShopUser(i IShopUser) {
	localShopUser = i
}

func ShopAdmin() IShopAdmin {
	if localShopAdmin == nil {
		panic("IShopAdmin interface not implemented yet")
	}
	return localShopAdmin
}

func InitShopAdmin(i IShopAdmin) {
	localShopAdmin = i
}

func ShopOwner() IShopOwner {
	if localShopOwner == nil {
		panic("IShopOwner interface not implemented yet")
	}
	return localShopOwner
}

func InitShopOwner(i IShopOwner) {
	localShopOwner = i
}

