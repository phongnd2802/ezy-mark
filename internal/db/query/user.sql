-- name: GetUserProfile :one
SELECT "user_id", "user_email", "user_nickname", "user_fullname", 
"user_avatar", "user_mobile", "user_gender", "user_birthday"
FROM "user_profile"
WHERE "user_id" = $1;

-- name: UpdateUserProfile :one
UPDATE "user_profile" 
SET 
    "user_nickname" = COALESCE(sqlc.arg(user_nickname), user_nickname),
    "user_fullname" = COALESCE(sqlc.arg(user_fullname), user_fullname),
    "user_avatar" = COALESCE(sqlc.arg(user_avatar), user_avatar),
    "user_mobile" = COALESCE(sqlc.arg(user_mobile), user_mobile),
    "user_gender" = COALESCE(sqlc.arg(user_gender), user_gender),
    "user_birthday" = COALESCE(sqlc.arg(user_birthday), user_birthday),
    "updated_at" = now()
WHERE "user_id" = sqlc.arg(user_id)
RETURNING *;


