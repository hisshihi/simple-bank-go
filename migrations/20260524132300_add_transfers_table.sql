-- +goose Up
CREATE TABLE transfers (
    id uuid primary key default uuidv7(),
    from_account_id uuid not null references accounts(id) on delete cascade,
    to_account_id uuid not null  references accounts(id) on delete cascade ,
    amount bigint not null ,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index idx_entries_from_account_id on transfers(from_account_id);
create index idx_entries_to_account_id on transfers(to_account_id);
create index idx_entries_account_ids on transfers(from_account_id, to_account_id);

-- +goose Down
SELECT 'down SQL query';
