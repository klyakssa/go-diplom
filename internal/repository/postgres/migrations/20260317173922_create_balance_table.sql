-- +goose Up
-- +goose StatementBegin
create table if not exists balances (
    id integer primary key generated always as identity,
    current INTEGER NOT NULL DEFAULT 0,
    user_id INTEGER REFERENCES users(id) NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists balances;
-- +goose StatementEnd
