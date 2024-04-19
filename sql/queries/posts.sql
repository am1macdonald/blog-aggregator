-- name: CreatePost :one
INSERT INTO posts 
(title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetPostsByUser :many
SELECT * FROM posts
  WHERE feed_id in (
    SELECT feed_id FROM feed_follows
    WHERE user_id = $1
  )
ORDER BY published_at; 
