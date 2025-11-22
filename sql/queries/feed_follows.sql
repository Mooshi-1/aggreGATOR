-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (
        id, created_at, updated_at, 
        user_id, feed_id
    ) VALUES (
        $1, NOW (), NOW (),
        $2, $3
) RETURNING *)
SELECT 
    inserted_feed_follow.*,
    feeds.name as feed_name,
    users.name as user_name
FROM inserted_feed_follow
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
INNER JOIN users ON inserted_feed_follow.user_id = users.id
;

-- name: GetFeedFollowsForUser :many
SELECT
    ff.id as ffid,
    f.name as feed_name,
    u.name as user_name
FROM feed_follows AS ff
INNER JOIN feeds AS f ON ff.feed_id = f.id
INNER JOIN users AS u ON ff.user_id = u.id
WHERE u.name = $1;

-- name: Unfollow :exec
DELETE FROM feed_follows
WHERE feed_id = $1
AND user_id = $2;