-- +goose Up

CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    title TEXT,
    url TEXT UNIQUE,
    description TEXT,
    published_at TIMESTAMPTZ,
    feed_id UUID
);

-- +goose Down
DROP TABLE posts;