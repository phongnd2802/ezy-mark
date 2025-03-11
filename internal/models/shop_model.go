package models

import "mime/multipart"

type RegisterShopReq struct {
	OwnerId         int64
	ShopName        string `json:"shop_name" form:"shop_name"`
	ShopDescription string `json:"shop_description" form:"shop_description"`
	ShopLogo        *multipart.FileHeader
	ShopPhone       string `json:"shop_phone" form:"shop_phone"`
	ShopEmail       string `json:"shop_email" form:"shop_email"`
	ShopAddress     string `json:"shop_address" form:"shop_address"`
	BusinessLicense *multipart.FileHeader
	TaxId           string `json:"tax_id" form:"tax_id"`
}

type ShopDetailsRes struct {
	ShopId          int64  `json:"shop_id"`
	ShopName        string `json:"shop_name"`
	ShopDescription string `json:"shop_description"`
	ShopLogo        string `json:"shop_logo"`
	ShopPhone       string `json:"shop_phone"`
	ShopEmail       string `json:"shop_email"`
	ShopAddress     string `json:"shop_address"`
}

type ApproveShopReq struct {
	ShopId int64 `json:"shop_id"`
	UserId int64
}

