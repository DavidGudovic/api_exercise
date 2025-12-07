-- +goose Up
-- +goose StatementBegin
ALTER TABLE tokens RENAME COLUMN expriry TO expiry;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tokens RENAME COLUMN expiry TO expriry;
-- +goose StatementEnd