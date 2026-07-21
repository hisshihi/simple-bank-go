-- +goose Up
CREATE TABLE entries
(
    id          uuid PRIMARY KEY     DEFAULT uuidv7(),
    account_id  uuid        not null references accounts (id) on delete cascade,
    transfer_id uuid        not null references  transfers(id) on delete cascade,
    amount      bigint      not null,
    currency    CHAR(3)     NOT NULL,
    direction   TEXT        NOT NULL CHECK ( direction IN ('debit', 'credit') ),
    created_at  timestamptz not null default now(),
    updated_at  timestamptz not null default now()
);

create index idx_entries_account_id on entries (account_id);
create index idx_entries_transfer_id on entries (transfer_id);

-- +goose Down
DROP TABLE entries;
