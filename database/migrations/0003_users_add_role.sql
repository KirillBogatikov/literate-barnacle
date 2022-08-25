-- +goose Up
alter table "users"
    add column "role" int not null default 1;

-- +goose Down
alter table "users"
    drop column "role";