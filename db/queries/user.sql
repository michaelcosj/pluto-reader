-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByOauthSub :one
SELECT * FROM users
WHERE oauth_sub = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  oauth_sub, email, name
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: CreateUserIfNotExist :exec
INSERT INTO users (
  oauth_sub, email, name
) VALUES (
  $1, $2, $3
) ON CONFLICT DO NOTHING;

-- name: UpdateUser :exec
UPDATE users
  set name = $2
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
