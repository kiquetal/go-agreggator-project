-- name: InsertFeedFollow :one

INSERT INTO follows_feeds (id,user_id, feed_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteFeedFollow :one

DELETE FROM follows_feeds
WHERE id = $1
RETURNING *;

-- name: GetAllFeedFollowsByUser :many

SELECT f.*
FROM feeds f
JOIN follows_feeds ff
ON f.id = ff.feed_id
WHERE ff.user_id = $1;
