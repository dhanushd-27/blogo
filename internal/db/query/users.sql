-- name: CreateUser :one
INSERT INTO users (
  name,
  email,
  password
)
VALUES ($1, $2, $3)
RETURNING id, name, email, password, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, name, email, password, created_at, updated_at
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, name, email, password, created_at, updated_at
FROM users
WHERE email = $1
LIMIT 1;

-- name: ListUsers :many
SELECT id, name, email, password, created_at, updated_at
FROM users
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET
  name = $2,
  email = $3,
  password = $4,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, email, password, created_at, updated_at;

-- name: DeleteUser :execrows
DELETE FROM users
WHERE id = $1;


