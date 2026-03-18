-- +goose Up
-- +goose StatementBegin
create table if not exists users (
    id INTEGER primary key GENERATED ALWAYS AS IDENTITY,
    login varchar(255) not null unique,
    password varchar(65) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
