-- name: CreateUser :one
INSERT INTO users(id,name,created_at,updated_at,api_key)
VALUES($1,$2,$3,$4,encode(sha256(random()::text::bytea), 'hex'))
RETURNING id,name,created_at,updated_at,api_key;

-- name: GetUserByApiKey :one
SELECT * FROM users
WHERE api_key = $1;
