-- name: CreateFeed :one  
INSERT INTO feeds (name, url, user_id)
VALUES ($1, $2, $3)
RETURNING *; 

-- name: CreateFeedFollow :one
INSERT INTO feed_follows (user_id, feed_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetAllFeeds :many
SELECT * FROM feeds;
