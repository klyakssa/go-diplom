-- +goose Up
-- +goose StatementBegin
ALTER TABLE withdraw_history
ALTER COLUMN sum TYPE NUMERIC(10, 2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE withdraw_history
ALTER COLUMN sum TYPE INTEGER;
-- +goose StatementEnd
