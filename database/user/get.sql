select u."id"       as "id",
       u."login"    as "login",
       u."password" as "password"
from "users" u
where login = $1
