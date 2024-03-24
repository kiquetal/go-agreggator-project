-- name: InsertFeed :one

INSERT INTO feeds (id,name,url,created_at, updated_at,user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllFeeds :many

SELECT * FROM feeds;


-- name: GetNexFeedsToFetch :many

SELECT * FROM feeds WHERE last_fetched_at IS NULL OR last_fetched_at < date_trunc('hour', now() - interval '1 week')
ORDER BY last_fetched_at ASC NULLS FIRST;

-- name: MarkedFetched :one

UPDATE feeds
SET last_fetched_at = $1 ,
   updated_at = $2
WHERE id = $3
RETURNING *;
