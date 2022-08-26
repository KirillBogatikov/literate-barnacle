select u."id"         as "id",
       u."login"      as "login",
       u."password"   as "password",
       u."role"       as "role",
       u."name"       as "name",
       u."surname"    as "surname",
       u."patronymic" as "patronymic",
       u."birth_date" as "birth_date"
from "users" u
where id = $1
