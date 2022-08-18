-- +goose Up
create table if not exists users
(
    id uuid not null primary key,
    login    text not null,
    password text not null
);

-- +goose Down
drop table users;