package user

type DbUser struct {
	Id       string `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
	Role     int    `db:"role"`
}
