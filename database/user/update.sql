update "users" u
set "login"      = coalesce(nullif(:login, ''), u."login"),
    "role"       = coalesce(nullif(:role, 0), u."role"),
    "name"       = coalesce(nullif(:name, ''), u."name"),
    "surname"    = coalesce(nullif(:surname, ''), u."surname"),
    "patronymic" = coalesce(nullif(:patronymic, ''), u."patronymic"),
    "birth_date" = coalesce(cast(nullif(:birth_date, '') as date), u."birth_date")
where u."id" = :id
returning
    u."id" as "id",
    u."login" as "login",
    u."password" as "password",
    u."role" as "role",
    u."name" as "name",
    u."surname" as "surname",
    u."patronymic" as "patronymic",
    u."birth_date" as "birth_date"