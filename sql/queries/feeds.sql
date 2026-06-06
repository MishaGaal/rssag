-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeedById :one
SELECT * FROM feeds WHERE id = $1;

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: GetUserFeeds :many
SELECT * FROM feeds WHERE user_id = $1;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
ORDER BY fetched_at ASC NULLS FIRST
LIMIT $1;

-- name: MarkFeedFetched :one
UPDATE feeds
SET fetched_at = NOW(), updated_at = NOW()
WHERE id = $1
RETURNING *;