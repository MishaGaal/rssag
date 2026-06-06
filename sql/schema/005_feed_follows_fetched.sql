-- +goose up
ALTER TABLE feeds ADD column fetched_at TIMESTAMP;

-- +goose down
ALTER TABLE feeds DROP column fetched_at;