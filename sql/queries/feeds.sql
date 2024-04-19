-- name: CreateFeed :one  
INSERT INTO feeds (name, url, user_id)
VALUES ($1, $2, $3)
RETURNING *; 

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: GetNextFeedsToFetch :many
SELECT *
FROM (SELECT *
      FROM feeds
      WHERE last_fetched_at IS NULL
      UNION
      SELECT *
      FROM feeds
      WHERE last_fetched_at IS NOT NULL
      ORDER BY last_fetched_at DESC) AS T1
LIMIT $1;

-- name: MarkFeedFetched :exec
UPDATE feeds SET
updated_at = NOW(),
last_fetched_at = NOW()
WHERE id = $1;

