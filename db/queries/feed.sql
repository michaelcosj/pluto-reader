-- name: FeedGetItems :many
SELECT * FROM feed_items
WHERE feed_id = $1;

-- name: FeedCreate :one
INSERT INTO feeds (
    title, description,
    site_link, feed_link
) VALUES (
    $1, $2, $3, $4
)
RETURNING id;

-- name: FeedAddItems :copyfrom
INSERT INTO feed_items (
    entry_id, title, summary, link,
    content, item_date, feed_id
) VALUES (
    $1, $2, $3, $4,
    $5, $6, $7 
);

