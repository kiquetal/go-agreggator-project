-- +goose UP
ALTER TABLE feeds ADD COLUMN test TEXT DEFAULT 'test';

-- +goose DOWN
ALTER TABLE feeds DROP COLUMN test;

