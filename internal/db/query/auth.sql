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

-- name: GetUserBaseById :one
SELECT *
FROM "user_base"
WHERE "user_id" = $1;

-- name: UpdateUserPassword :exec
UPDATE "user_base"
SET "user_password" = $1, "updated_at" = now()
WHERE "user_id" = $2;

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


-- name: CreateUserSession :one
INSERT INTO "user_session" (
    "sub_token",
    "refresh_token",
    "user_agent",
    "client_ip",
    "user_login_time",
    "expires_at",
    "user_id"
) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetUserByUserHash :one
SELECT *
FROM "user_base"
WHERE "user_hash" = $1;


-- name: DeleteSessionBySubToken :exec
DELETE FROM "user_session"
WHERE "sub_token" = $1;


-- name: DeleteSessionByUserId :exec
DELETE FROM "user_session"
WHERE "user_id" = $1;

-- name: CheckRefreshTokenUsed :one
SELECT COUNT(*)
FROM "user_session"
WHERE "refresh_token_used" = $1;

-- name: GetSessionBySubToken :one
SELECT "session_id", "user_id", "refresh_token", "refresh_token_used"
FROM "user_session"
WHERE "sub_token" = $1 
LIMIT 1;

-- name: GetSessionByRefreshTokenUsed :one
SELECT "session_id", "user_id", "refresh_token", "refresh_token_used"
FROM "user_session"
WHERE "refresh_token_used" = $1;

-- name: UpdateSession :exec
UPDATE "user_session"
SET "refresh_token" = $1, "refresh_token_used" = $2, "expires_at" = $3, "sub_token" = $4
WHERE "session_id" = $5;


