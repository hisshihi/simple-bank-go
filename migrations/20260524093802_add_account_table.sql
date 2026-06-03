-- +goose Up
CREATE TABLE accounts
(
    id         UUID PRIMARY KEY        DEFAULT uuidv7(),
    owner      TEXT           NOT NULL,
    balance    BIGINT NOT NULL CHECK ( balance >= 0 ),
    currency   CHAR(3)        NOT NULL,
    status     TEXT           NOT NULL,
    created_at timestamptz    NOT NULL default now(),
    updated_at timestamptz    NOT NULL default now(),
    deleted_at timestamptz
);

CREATE INDEX idx_accounts_status ON accounts (status);

-- +goose Down
DROP TABLE accounts;
