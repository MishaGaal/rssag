-- +goose up
CREATE TABLE feeds (
    Id  UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users(Id) ON DELETE CASCADE

);

-- +goose down
 DROP TABLE  feeds;