-- +goose Up
-- +goose StatementBegin
create table if not exists withdraw_history (
    number integer primary key,
    sum INTEGER NOT NULL,
    user_id INTEGER REFERENCES users(id) NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists withdraw_history;
-- +goose StatementEnd
