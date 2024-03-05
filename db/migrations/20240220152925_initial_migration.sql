-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id          INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    oauth_sub   VARCHAR(255) UNIQUE NOT NULL,
    email       VARCHAR(255) UNIQUE NOT NULL,
    name        VARCHAR(255) NOT NULL
);

CREATE TABLE feeds (
    id              INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title           VARCHAR(255),
    description     TEXT,
    site_link       VARCHAR(255),
    feed_link       VARCHAR(255) UNIQUE NOT NULL,
    last_refreshed  TIMESTAMPTZ 
);

CREATE TABLE feed_items (
    id              INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    entry_id        VARCHAR(255),
    title           VARCHAR(255),
    summary         VARCHAR(255),
    link            VARCHAR(255) UNIQUE NOT NULL,
    content         TEXT,
    item_date       TIMESTAMPTZ,
    feed_id         INTEGER NOT NULL,

    FOREIGN KEY(feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feed_items;
DROP TABLE feeds;
DROP TABLE users;
-- +goose StatementEnd
