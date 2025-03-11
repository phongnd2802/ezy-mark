-- name: CheckShopExist :one
SELECT COUNT(*)
FROM "shops"
WHERE "shop_email" = $1 OR "tax_id" = $2;


-- name: CreateShop :one
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
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;



-- name: CheckShopExistByAdmin :one
SELECT COUNT(*)
FROM "shops"
WHERE "shop_id" = $1;

-- name: ApproveShop :exec
UPDATE "shops"
SET "is_verified" = true,
    "verified_by" = $1,
    "verified_at" = now()
WHERE "shop_id" = $2;
