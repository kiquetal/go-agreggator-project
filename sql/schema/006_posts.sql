-- +goose Up
CREATE TABLE posts (
    id uuid PRIMARY KEY ,
    created_at timestamp DEFAULT now(),
    updated_at timestamp DEFAULT now(),
    title text,
    url text,
    description text,
    published_at timestamp,
    feed_id uuid REFERENCES feeds(id)
);

-- +goose Down

DROP TABLE posts;
