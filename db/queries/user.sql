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

-- name: UserGetFeedItems :many
SELECT
	fi.id, fi.title, fi.summary,
    fi.link, fi.item_date, fi.content,
    ufi.is_read, uf.feed_name
FROM
	feed_items fi
JOIN user_feeds uf on 
    (uf.feed_id = fi.feed_id)
JOIN user_feed_items ufi on
	(ufi.item_id = fi.id)
WHERE
	ufi.user_id = $1;

-- name: UserGetFeedItemContent :one
SELECT fi.content FROM
    feed_items fi
JOIN user_feed_items ufi on
	(ufi.user_id = $1)
WHERE 
    fi.id = $2;
