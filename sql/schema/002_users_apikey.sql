-- +goose up
ALTER TABLE users ADD column api_key VARCHAR(64) UNIQUE NOT NULL DEFAULT (
    encode(sha256(random()::text::bytea), 'hex')
    );

-- +goose down
ALTER TABLE users DROP column api_key;