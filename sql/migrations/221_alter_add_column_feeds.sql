-- +goose UP
ALTER TABLE feeds ADD COLUMN test2 TEXT DEFAULT 'test';

-- +goose DOWN
ALTER TABLE feeds DROP COLUMN test2;

