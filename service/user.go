package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	jwt "github.com/Viva-Victoria/bear-jwt"
	"github.com/Viva-Victoria/bear-jwt/alg"
	"github.com/google/uuid"

	"literate-barnacle/database"
	"literate-barnacle/database/user"
	"literate-barnacle/service/hash"
)

type UserService interface {
	Login(ctx context.Context, request LoginRequest) (LoginResponse, error)
	Signup(ctx context.Context, request SignUpRequest) (SignUpResponse, error)
}

type UserServiceImpl struct {
	repo      user.Repository
	encryptor hash.Encryptor
}

func NewUserServiceImpl(repo user.Repository, encryptor hash.Encryptor) UserServiceImpl {
	return UserServiceImpl{
		repo:      repo,
		encryptor: encryptor,
	}
}

func (u UserServiceImpl) Login(ctx context.Context, request LoginRequest) (LoginResponse, error) {
	if validation := request.Validate(); !validation.IsValid() {
		return LoginResponse{
			BaseResponse: BaseResponse{
				Error:      "Введены некорректные данные",
				Validation: validation,
			},
		}, nil
	}

	dbUser, err := u.repo.Get(ctx, request.Login)
	switch {
	case errors.Is(err, database.ErrNotFound):
		return LoginResponse{
			BaseResponse: BaseResponse{
				Error: fmt.Sprintf("Пользователь \"%s\" не найден", request.Login),
			},
		}, nil
	case err != nil:
		return LoginResponse{}, fmt.Errorf("can't get user from DB: %v", err)
	}

	err = u.encryptor.Compare(dbUser.Password, request.Password)
	switch {
	case errors.Is(err, hash.ErrMismatched):
		return LoginResponse{
			BaseResponse: BaseResponse{
				Error: "Неверный пароль",
			},
		}, nil
	case err != nil:
		return LoginResponse{}, fmt.Errorf("can't compare passwords: %v", err)
	}

	now := time.Now()
	token := jwt.NewToken(alg.EdDSA)

	claims := TokenClaims{
		BasicClaims: jwt.BasicClaims{
			Id:        uuid.NewString(),
			IssuedAt:  jwt.NewPosixTime(now),
			NotBefore: jwt.NewPosixTime(now.Add(time.Second * 15)),
			ExpiresAt: jwt.NewPosixTime(now.Add(time.Minute * 30)),
		},
		UserId: dbUser.Id,
	}

	err = token.Claims.Set(claims)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("can't create token: %v", err)
	}

	tokenString, err := token.WriteString()
	if err != nil {
		return LoginResponse{}, fmt.Errorf("can't create token: %v", err)
	}

	return LoginResponse{
		Token: tokenString,
	}, nil
}

func (u UserServiceImpl) Signup(ctx context.Context, request SignUpRequest) (SignUpResponse, error) {
	if validation := request.Validate(); !validation.IsValid() {
		return SignUpResponse{
			BaseResponse: BaseResponse{
				Error:      "Введены некорректные данные",
				Validation: validation,
			},
		}, nil
	}

	existUser, err := u.repo.Get(ctx, request.User.Credentials.Login)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return SignUpResponse{}, fmt.Errorf("can't check login: %v", err)
	}
	if len(existUser.Id) > 0 {
		return SignUpResponse{
			BaseResponse: BaseResponse{
				Error: "Данный логин занят",
			},
		}, nil
	}

	passwordHash, err := u.encryptor.Encrypt(request.User.Credentials.Password)
	if err != nil {
		return SignUpResponse{}, err
	}

	id := uuid.New()
	err = u.repo.Add(ctx, user.DbUser{
		Id:       id.String(),
		Login:    request.User.Credentials.Login,
		Password: passwordHash,
	})
	if err != nil {
		return SignUpResponse{}, err
	}

	return SignUpResponse{
		UserId: id,
	}, nil
}
