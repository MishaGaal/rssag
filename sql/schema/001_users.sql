-- +goose up
CREATE TABLE users (
    Id  UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL
);

-- +goose down

 DROP TABLE  users