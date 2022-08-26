package user

type DbUser struct {
	Id         string `db:"id,omitempty"`
	Login      string `db:"login,omitempty"`
	Password   string `db:"password,omitempty"`
	Name       string `db:"name,omitempty"`
	Surname    string `db:"surname,omitempty"`
	Patronymic string `db:"patronymic,omitempty"`
	BirthDate  string `db:"birth_date,omitempty"`
	Role       uint8  `db:"role,omitempty"`
}
