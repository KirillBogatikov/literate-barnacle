package user

import (
	"errors"
	"fmt"
	"literate-barnacle/service"
	"literate-barnacle/service/ctx"
	"literate-barnacle/service/models"
	"time"

	jwt "github.com/Viva-Victoria/bear-jwt"
	"github.com/Viva-Victoria/bear-jwt/alg"
	"github.com/google/uuid"

	"literate-barnacle/database"
	"literate-barnacle/database/user"
	"literate-barnacle/service/hash"
)

type Service interface {
	Login(ctx ctx.Context, request LoginRequest) (LoginResponse, error)
	Signup(ctx ctx.Context, request SignUpRequest) (SignUpResponse, error)
	Get(ctx ctx.Context, id uuid.UUID) (Response, error)
	Update(ctx ctx.Context, user models.User) (Response, error)
}

type ServiceImpl struct {
	repo      user.Repository
	encryptor hash.Encryptor
}

func NewServiceImpl(repo user.Repository, encryptor hash.Encryptor) ServiceImpl {
	return ServiceImpl{
		repo:      repo,
		encryptor: encryptor,
	}
}

func (u ServiceImpl) Login(ctx ctx.Context, request LoginRequest) (LoginResponse, error) {
	if validation := request.Validate(); !validation.IsValid() {
		return LoginResponse{
			BaseResponse: service.BaseResponse{
				Error:      "Введены некорректные данные",
				Validation: validation,
			},
		}, nil
	}

	dbUser, err := u.repo.Get(ctx, request.Login)
	switch {
	case errors.Is(err, database.ErrNotFound):
		return LoginResponse{
			BaseResponse: service.BaseResponse{
				Error: fmt.Sprintf("Пользователь \"%s\" не найден", request.Login),
			},
		}, nil
	case err != nil:
		return LoginResponse{}, fmt.Errorf("can't get user from DB: %w", err)
	}

	err = u.encryptor.Compare(dbUser.Password, request.Password)
	switch {
	case errors.Is(err, hash.ErrMismatched):
		return LoginResponse{
			BaseResponse: service.BaseResponse{
				Error: "Неверный пароль",
			},
		}, nil
	case err != nil:
		return LoginResponse{}, fmt.Errorf("can't compare passwords: %w", err)
	}

	now := time.Now()
	token := jwt.NewToken(alg.EdDSA)

	userId, err := uuid.Parse(dbUser.Id)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("can't parse user id: %w", err)
	}

	claims := models.TokenClaims{
		BasicClaims: jwt.BasicClaims{
			Id:        uuid.NewString(),
			IssuedAt:  jwt.NewPosixTime(now),
			NotBefore: jwt.NewPosixTime(now.Add(time.Second * 15)),
			ExpiresAt: jwt.NewPosixTime(now.Add(time.Minute * 30)),
		},
		Authorization: models.Authorization{
			UserId: userId,
			Role:   models.Role(dbUser.Role),
		},
	}

	err = token.Claims.Set(claims)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("can't create token: %w", err)
	}

	tokenString, err := token.WriteString()
	if err != nil {
		return LoginResponse{}, fmt.Errorf("can't create token: %w", err)
	}

	return LoginResponse{
		Token: tokenString,
	}, nil
}

func (u ServiceImpl) Signup(ctx ctx.Context, request SignUpRequest) (SignUpResponse, error) {
	if validation := request.Validate(); !validation.IsValid() {
		return SignUpResponse{
			BaseResponse: service.BaseResponse{
				Error:      "Введены некорректные данные",
				Validation: validation,
			},
		}, nil
	}

	existUser, err := u.repo.Get(ctx, request.User.Credentials.Login)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return SignUpResponse{}, fmt.Errorf("can't check login: %w", err)
	}
	if len(existUser.Id) > 0 {
		return SignUpResponse{
			BaseResponse: service.BaseResponse{
				Error: "Данный логин занят",
			},
		}, nil
	}

	passwordHash, err := u.encryptor.Encrypt(request.User.Credentials.Password)
	if err != nil {
		return SignUpResponse{}, fmt.Errorf("can't get hash: %w", err)
	}

	id := uuid.New()
	err = u.repo.Add(ctx, user.DbUser{
		Id:       id.String(),
		Login:    request.User.Credentials.Login,
		Password: passwordHash,
	})
	if err != nil {
		return SignUpResponse{}, fmt.Errorf("can't save user: %w", err)
	}

	return SignUpResponse{
		UserId: id,
	}, nil
}

func (u ServiceImpl) Get(c ctx.Context, id uuid.UUID) (Response, error) {
	if !c.Authorized {
		return Response{
			BaseResponse: service.BaseResponse{
				Error: "Необходимо выполнить вход",
			},
		}, ctx.ErrUnauthorized
	}
	if c.Authorization.UserId.String() != id.String() && c.Authorization.Role != models.RoleAdmin {
		return Response{
			BaseResponse: service.BaseResponse{
				Error: "Отказано в доступе",
			},
		}, ctx.ErrForbidden
	}

	dbUser, err := u.repo.GetById(c, id.String())
	switch {
	case errors.Is(err, database.ErrNotFound):
		return Response{
			BaseResponse: service.BaseResponse{
				Error: "Пользователь не найден",
			},
		}, database.ErrNotFound
	case err != nil:
		return Response{}, fmt.Errorf("can't get user: %w", err)
	}

	domainUser, err := mapDbUser(dbUser)
	if err != nil {
		return Response{}, fmt.Errorf("can't map user: %w", err)
	}

	domainUser.Credentials.Password = ""
	return Response{
		User: &domainUser,
	}, nil
}

func (u ServiceImpl) Update(c ctx.Context, user models.User) (Response, error) {
	if !c.Authorized {
		return Response{
			BaseResponse: service.BaseResponse{
				Error: "Необходимо выполнить вход",
			},
		}, ctx.ErrUnauthorized
	}
	if c.Authorization.UserId.String() != user.Id.String() && c.Authorization.Role != models.RoleAdmin {
		return Response{
			BaseResponse: service.BaseResponse{
				Error: "Отказано в доступе",
			},
		}, ctx.ErrForbidden
	}

	dbUser, err := u.repo.Update(c, mapUser(user))
	switch {
	case errors.Is(err, database.ErrNotFound):
		return Response{
			BaseResponse: service.BaseResponse{
				Error: "Пользователь не найден",
			},
		}, database.ErrNotFound
	case err != nil:
		return Response{}, fmt.Errorf("can't update user: %w", err)
	}

	domainUser, err := mapDbUser(dbUser)
	if err != nil {
		return Response{}, fmt.Errorf("can't map user: %w", err)
	}

	domainUser.Credentials.Password = ""
	return Response{
		User: &domainUser,
	}, nil
}
