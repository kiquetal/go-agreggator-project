-- +goose UP
ALTER table feeds
ADD COLUMN last_fetched_at TIMESTAMP WITH TIME ZONE;

-- +goose DOWN
ALTER table feeds
DROP COLUMN last_fetched_at;
