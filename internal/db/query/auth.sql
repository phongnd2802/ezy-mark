-- name: CheckUserBaseExists :one
SELECT COUNT(*)
FROM "user_base"
WHERE "user_email" = $1;

-- name: CreateUserBase :one
INSERT INTO "user_base" (
    "user_email",
    "user_password",
    "user_hash",
    "user_otp"
) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetUserBaseByEmail :one
SELECT * 
FROM "user_base"
WHERE "user_email" = $1;

-- name: UpdateUserVerify :one
UPDATE "user_base"
SET "is_verified" = true
WHERE "user_hash" = $1 RETURNING *;

-- name: CreateUserProfile :one
INSERT INTO "user_profile" (
    "user_id",
    "user_email",
    "user_nickname"
) VALUES ($1, $2, $3) RETURNING *;
