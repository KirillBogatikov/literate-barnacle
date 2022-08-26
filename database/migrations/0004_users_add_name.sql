-- +goose Up
alter table "users"
    add column "name" text not null default '',
    add column "surname" text not null default '',
    add column "patronymic" text not null default '',
    add column "birth_date" date not null default now();

-- +goose Down
alter table "users"
    drop column "name",
    drop column "surname",
    drop column "patronymic",
    drop column "birth_date";