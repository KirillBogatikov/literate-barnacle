-- +goose Up
create unique index if not exists "users_login" on "users"("login");

-- +goose Down
drop index if exists "users_login";