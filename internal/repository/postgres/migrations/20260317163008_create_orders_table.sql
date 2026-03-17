-- +goose Up
-- +goose StatementBegin
create table if not exists orders (
    number VARCHAR(1024) PRIMARY KEY,
    status VARCHAR(15) NOT NULL DEFAULT 'NEW',
    accrual INTEGER DEFAULT 0,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id INTEGER REFERENCES users(id) NOT NULL,

    CONSTRAINT status_check CHECK (
        status IN ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED')
    )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists orders;
-- +goose StatementEnd
