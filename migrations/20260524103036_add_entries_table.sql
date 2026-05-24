-- +goose Up
CREATE TABLE entries (
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    account_id uuid not null references accounts(id) on delete cascade ,
    amount bigint not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index idx_entries_account_id on entries(account_id);

-- +goose Down
DROP TABLE entries;
