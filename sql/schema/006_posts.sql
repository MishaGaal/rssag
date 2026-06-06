-- +goose up
CREATE TABLE posts (
    Id  UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    published_at TIMESTAMP,
    url TEXT UNIQUE NOT NULL UNIQUE,
    feed_id UUID NOT NULL REFERENCES feeds(Id) ON DELETE CASCADE
);

-- +goose down
 DROP TABLE posts;