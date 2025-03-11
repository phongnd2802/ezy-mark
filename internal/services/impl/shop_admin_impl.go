package impl

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/phongnd2802/ezy-mark/internal/db"
	"github.com/phongnd2802/ezy-mark/internal/models"
	"github.com/phongnd2802/ezy-mark/internal/response"
	"github.com/phongnd2802/ezy-mark/internal/services"
)

type sShopAdmin struct {
	store db.Store
}

// ApproveShop implements services.IShopAdmin.
func (s *sShopAdmin) ApproveShop(ctx context.Context, params *models.ApproveShopReq) (int, error) {
	exist, err := s.store.CheckShopExistByAdmin(ctx, params.ShopId)
	if err != nil {
		return response.ErrCodeInternalServer, err
	}
	if exist == 0 {
		return response.ErrCodeBadRequest, fmt.Errorf("shop with ID %d does not exist", params.ShopId)
	}

	err = s.store.ApproveShop(ctx, db.ApproveShopParams{
		VerifiedBy: pgtype.Int8{Int64: params.UserId, Valid: true},
		ShopID: params.ShopId,
	})
	if err != nil {
		return response.ErrCodeInternalServer, err
	}

	return response.ErrCodeSuccess, nil
}

// BlockShop implements services.IShopAdmin.
func (s *sShopAdmin) BlockShop(ctx context.Context, shopId int64) (code int, err error) {
	panic("unimplemented")
}

// GetAllShops implements services.IShopAdmin.
func (s *sShopAdmin) GetAllShops(ctx context.Context) (code int, res []db.Shop, err error) {
	panic("unimplemented")
}

// GetShopById implements services.IShopAdmin.
func (s *sShopAdmin) GetShopById(ctx context.Context, shopId int64) (code int, res db.Shop, err error) {
	panic("unimplemented")
}

func NewShopAdminImpl(store db.Store) *sShopAdmin {
	return &sShopAdmin{
		store: store,
	}
}

var _ services.IShopAdmin = (*sShopAdmin)(nil)
