package user

import (
	"literate-barnacle/database/user"
	"literate-barnacle/service/models"
	"time"

	"github.com/google/uuid"
)

var (
	_dateLayout = "2006-01-02"
)

func mapDbUser(user user.DbUser) (models.User, error) {
	id, err := uuid.Parse(user.Id)
	if err != nil {
		return models.User{}, err
	}

	birthDate, err := time.Parse(time.RFC3339, user.BirthDate)
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		Id: id,
		Credentials: models.Credentials{
			Login:    user.Login,
			Password: user.Password,
		},
		Role:       models.Role(user.Role),
		Name:       user.Name,
		Surname:    user.Surname,
		Patronymic: user.Patronymic,
		BirthDate:  birthDate,
	}, nil
}

func mapUser(u models.User) user.DbUser {
	return user.DbUser{
		Id:         u.Id.String(),
		Login:      u.Credentials.Login,
		Password:   u.Credentials.Password,
		Name:       u.Name,
		Surname:    u.Surname,
		Patronymic: u.Patronymic,
		BirthDate:  u.BirthDate.Format(_dateLayout),
		Role:       uint8(u.Role),
	}
}
