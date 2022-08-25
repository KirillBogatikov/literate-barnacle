package user

import (
	"literate-barnacle/database/user"
	"literate-barnacle/service/models"

	"github.com/google/uuid"
)

func mapDbUser(user user.DbUser) (models.User, error) {
	id, err := uuid.Parse(user.Id)
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		Id: id,
		Credentials: models.Credentials{
			Login: user.Login,
		},
		Role: models.Role(user.Role),
	}, nil
}
