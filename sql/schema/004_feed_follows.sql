-- +goose up
CREATE TABLE feed_follows (
    Id  UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(Id) ON DELETE CASCADE,
    feed_id UUID NOT NULL REFERENCES feeds(Id) ON DELETE CASCADE,
    UNIQUE(user_id, feed_id)
);

-- +goose down
 DROP TABLE  feed_follows;