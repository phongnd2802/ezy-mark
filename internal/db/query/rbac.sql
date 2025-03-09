-- name: GetRoleByUserId :many
SELECT r.role_name
FROM "user_roles" ur
JOIN "roles" r ON ur.role_id = r.role_id
WHERE ur.user_id = $1;
