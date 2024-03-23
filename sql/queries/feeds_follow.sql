-- name: InsertFeedFollow :one

INSERT INTO follows_feeds (id,user_id, feed_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
