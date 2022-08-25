select u."id"       as "id",
       u."login"    as "login",
       u."password" as "password",
       u."role"     as "role"
from "users" u
where id = $1
