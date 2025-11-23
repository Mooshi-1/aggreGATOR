-- name: CreatePost :exec
INSERT INTO posts (
    id,
    created_at,
    updated_at,
    title,
    url,
    description,
    published_at,
    feed_id

) VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    $4,
    $5,
    $6
);

-- name: GetPostsForUser :many
SELECT p.id, p.title, p.url, p.feed_id, f.user_id
FROM posts AS p
INNER JOIN feeds AS f ON p.feed_id = f.feed_id
WHERE user_id = $1;