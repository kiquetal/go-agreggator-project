-- +goose Up

CREATE TABLE follows_feeds(
                      id UUID PRIMARY KEY,
                      created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                      updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                      user_id UUID NOT NULL,
                      feed_id UUID NOT NULL,
                      FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                      FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE

);

-- +goose Down

DROP TABLE follows_feeds;


