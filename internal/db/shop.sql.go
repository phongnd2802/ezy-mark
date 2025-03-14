// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: shop.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const approveShop = `-- name: ApproveShop :exec
UPDATE "shops"
SET "is_verified" = true,
    "verified_by" = $1,
    "verified_at" = now()
WHERE "shop_id" = $2
`

type ApproveShopParams struct {
	VerifiedBy pgtype.Int8 `json:"verified_by"`
	ShopID     int64       `json:"shop_id"`
}

func (q *Queries) ApproveShop(ctx context.Context, arg ApproveShopParams) error {
	_, err := q.db.Exec(ctx, approveShop, arg.VerifiedBy, arg.ShopID)
	return err
}

const checkShopExist = `-- name: CheckShopExist :one
SELECT COUNT(*)
FROM "shops"
WHERE "shop_email" = $1 OR "tax_id" = $2
`

type CheckShopExistParams struct {
	ShopEmail string `json:"shop_email"`
	TaxID     string `json:"tax_id"`
}

func (q *Queries) CheckShopExist(ctx context.Context, arg CheckShopExistParams) (int64, error) {
	row := q.db.QueryRow(ctx, checkShopExist, arg.ShopEmail, arg.TaxID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const checkShopExistByAdmin = `-- name: CheckShopExistByAdmin :one
SELECT COUNT(*)
FROM "shops"
WHERE "shop_id" = $1
`

func (q *Queries) CheckShopExistByAdmin(ctx context.Context, shopID int64) (int64, error) {
	row := q.db.QueryRow(ctx, checkShopExistByAdmin, shopID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createShop = `-- name: CreateShop :one
INSERT INTO "shops" (
    "owner_id",
    "shop_name",
    "shop_description",
    "shop_logo",
    "shop_phone",
    "shop_email",
    "shop_address",
    "business_license",
    "tax_id"
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING shop_id, owner_id, shop_name, shop_description, shop_logo, shop_phone, shop_email, shop_address, business_license, tax_id, is_active, is_blocked, is_verified, verified_by, verified_at, created_at, updated_at
`

type CreateShopParams struct {
	OwnerID         int64       `json:"owner_id"`
	ShopName        string      `json:"shop_name"`
	ShopDescription pgtype.Text `json:"shop_description"`
	ShopLogo        string      `json:"shop_logo"`
	ShopPhone       pgtype.Text `json:"shop_phone"`
	ShopEmail       string      `json:"shop_email"`
	ShopAddress     string      `json:"shop_address"`
	BusinessLicense string      `json:"business_license"`
	TaxID           string      `json:"tax_id"`
}

func (q *Queries) CreateShop(ctx context.Context, arg CreateShopParams) (Shop, error) {
	row := q.db.QueryRow(ctx, createShop,
		arg.OwnerID,
		arg.ShopName,
		arg.ShopDescription,
		arg.ShopLogo,
		arg.ShopPhone,
		arg.ShopEmail,
		arg.ShopAddress,
		arg.BusinessLicense,
		arg.TaxID,
	)
	var i Shop
	err := row.Scan(
		&i.ShopID,
		&i.OwnerID,
		&i.ShopName,
		&i.ShopDescription,
		&i.ShopLogo,
		&i.ShopPhone,
		&i.ShopEmail,
		&i.ShopAddress,
		&i.BusinessLicense,
		&i.TaxID,
		&i.IsActive,
		&i.IsBlocked,
		&i.IsVerified,
		&i.VerifiedBy,
		&i.VerifiedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
