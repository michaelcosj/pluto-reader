-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_feeds (
    id              INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id         INTEGER NOT NULL,
    feed_id         INTEGER NOT NULL,
    feed_name       VARCHAR(255),
    last_updated    TIMESTAMPTZ,
    update_interval INTERVAL NOT NULL,

    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(feed_id) REFERENCES feeds(id),

    UNIQUE (user_id, feed_id)
);

CREATE TABLE user_feed_items (
    id              INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id         INTEGER NOT NULL,
    item_id         INTEGER NOT NULL,
    is_read         BOOLEAN DEFAULT false,

    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(item_id) REFERENCES feed_items(id),

    UNIQUE (user_id, item_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_feeds;
DROP TABLE user_feed_items;
-- +goose StatementEnd
