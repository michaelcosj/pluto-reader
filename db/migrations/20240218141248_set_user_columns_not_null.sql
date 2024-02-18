-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
     ALTER COLUMN email SET NOT NULL;

ALTER TABLE users
     ALTER COLUMN name SET NOT NULL;

ALTER TABLE users
     ALTER COLUMN oauth_sub SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
     ALTER COLUMN email DROP NOT NULL;

ALTER TABLE users
     ALTER COLUMN name DROP NOT NULL;

ALTER TABLE users
     ALTER COLUMN oauth_sub DROP NOT NULL;
-- +goose StatementEnd
