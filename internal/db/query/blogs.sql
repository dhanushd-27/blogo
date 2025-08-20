-- name: CreateBlog :one
INSERT INTO blogs (
  title,
  content,
  user_id
)
VALUES ($1, $2, $3)
RETURNING id, title, content, user_id, created_at, updated_at;

-- name: GetBlogByID :one
SELECT id, title, content, user_id, created_at, updated_at
FROM blogs
WHERE id = $1
LIMIT 1;

-- name: ListBlogs :many
SELECT id, title, content, user_id, created_at, updated_at
FROM blogs
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: ListBlogsByUserID :many
SELECT id, title, content, user_id, created_at, updated_at
FROM blogs
WHERE user_id = $1
ORDER BY id
LIMIT $2 OFFSET $3;

-- name: UpdateBlog :one
UPDATE blogs
SET
  title = $2,
  content = $3,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, title, content, user_id, created_at, updated_at;

-- name: DeleteBlog :execrows
DELETE FROM blogs
WHERE id = $1;


