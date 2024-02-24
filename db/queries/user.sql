-- name: UserGetByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: UserGetByOauthSub :one
SELECT * FROM users
WHERE oauth_sub = $1 LIMIT 1;

-- name: UserCreate :one
INSERT INTO users (
  oauth_sub, email, name
) VALUES (
  $1, $2, $3
)
RETURNING id;

-- name: UserUpdate :exec
UPDATE users
  set name = $2
WHERE id = $1;

-- name: UserDelete :exec
DELETE FROM users
WHERE id = $1;

-- name: UserAddFeed :exec
INSERT INTO user_feeds (
    user_id, feed_id,
    feed_name, update_interval
) VALUES (
    $1, $2, $3, $4
);


-- name: UserAddFeedItems :copyfrom
INSERT INTO user_feed_items (
    user_id, item_id
) VALUES (
    $1, $2
);

