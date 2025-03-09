package impl

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minio/minio-go/v7"
	"github.com/phongnd2802/ezy-mark/internal/consts"
	"github.com/phongnd2802/ezy-mark/internal/db"
	"github.com/phongnd2802/ezy-mark/internal/helpers"
	"github.com/phongnd2802/ezy-mark/internal/models"
	"github.com/phongnd2802/ezy-mark/internal/pkg/utils"
	"github.com/phongnd2802/ezy-mark/internal/response"
	"github.com/phongnd2802/ezy-mark/internal/services"
)

type sShopUser struct {
	store db.Store
}

// GetShopDetails implements services.IShopUser.
func (s *sShopUser) GetShopDetails(ctx context.Context, shopId int64) (code int, res *models.ShopDetailsRes, err error) {
	panic("unimplemented")
}

// RegisterShop implements services.IShopUser.
func (s *sShopUser) RegisterShop(ctx context.Context, params *models.RegisterShopReq) (int, error) {
	foundShop, err := s.store.CheckShopExist(ctx, db.CheckShopExistParams{
		ShopEmail: params.ShopEmail,
		TaxID:    params.TaxId,
	})
	if err != nil {
		return response.ErrCodeInternalServer, err
	}
	if foundShop > 0 {
		return response.ErrCodeShopEmailAlreadyExists, nil
	}

	shopLogo := helpers.GenerateShopLogoObjectName(params.ShopName, params.ShopLogo.Filename)
	businessLicense := helpers.GenerateBusinessLicenseObjectName(params.ShopName, params.BusinessLicense.Filename)

	_, err = s.store.CreateShop(ctx, db.CreateShopParams{
		OwnerID: params.OwnerId,
		ShopName: params.ShopName,
		ShopEmail: params.ShopEmail,
		ShopPhone: pgtype.Text{String: params.ShopPhone, Valid: params.ShopPhone != ""},
		ShopAddress: params.ShopAddress,
		ShopDescription: pgtype.Text{String: params.ShopDescription, Valid: params.ShopDescription != ""},
		TaxID: params.TaxId,
		ShopLogo: utils.GetURLFileName(consts.BucketShopLogo, shopLogo),
		BusinessLicense: utils.GetURLFileName(consts.BucketShopBusinessLicense, businessLicense),
	})	
	if err != nil {
		return response.ErrCodeInternalServer, err
	}

	// Upload Shop Logo to Minio
	err = services.UploadService().UploadFile(ctx, consts.BucketShopLogo, shopLogo, params.ShopLogo, minio.PutObjectOptions{
		ContentType: params.ShopLogo.Header.Get("Content-Type"),
	})
	if err != nil {
		return response.ErrCodeInternalServer, err
	}

	// Upload Business License to Minio
	err = services.UploadService().UploadFile(ctx, consts.BucketShopBusinessLicense, businessLicense, params.BusinessLicense, minio.PutObjectOptions{
		ContentType: params.BusinessLicense.Header.Get("Content-Type"),
	})
	if err != nil {
		return response.ErrCodeInternalServer, err
	}

	return response.ErrCodeSuccess, nil
}

func NewShopUser(store db.Store) *sShopUser {
	return &sShopUser{
		store: store,
	}
}

var _ services.IShopUser = (*sShopUser)(nil)
