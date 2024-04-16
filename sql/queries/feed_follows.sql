-- name: CreateFeedFollow :one
INSERT INTO feed_follows (user_id, feed_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetAllFeedFollows :many
SELECT * FROM feed_follows
WHERE user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE feed_id = $1
AND user_id = $2;

