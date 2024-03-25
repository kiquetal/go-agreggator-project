-- name: InsertPost :one

INSERT INTO posts (id,title,url,description,published_at,feed_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;


-- name: GetPostByUsers :many

SELECT p.*
FROM posts p
JOIN feeds f ON p.feed_id = f.id
JOIN follows_feeds uf ON f.id = uf.feed_id
WHERE uf.user_id = $1
ORDER BY p.published_at DESC;

