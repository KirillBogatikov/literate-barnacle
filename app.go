package main

import (
	"encoding/base64"
	"fmt"
	"literate-barnacle/api"
	"literate-barnacle/config"
	"literate-barnacle/database"
	"literate-barnacle/database/user"
	"literate-barnacle/service"
	"literate-barnacle/service/hash"
	"time"

	jwt "github.com/Viva-Victoria/bear-jwt"
	"github.com/Viva-Victoria/bear-jwt/alg"
	"github.com/jmoiron/sqlx"
)

type App struct {
	settings       config.Settings
	timeoutContext service.TimeoutContextProvider

	postgres       *sqlx.DB
	userRepository user.Repository

	encryptor   hash.Encryptor
	userService service.UserService

	server api.Server
}

func NewApp(settings config.Settings, timeoutContextProvider service.TimeoutContextProvider) *App {
	return &App{
		settings:       settings,
		timeoutContext: timeoutContextProvider,
	}
}

func (a *App) InitRepositories() error {
	postgresCtx, cancelPostgresCtx := a.timeoutContext(time.Second * 15)
	defer cancelPostgresCtx()

	postgres, err := database.NewPgx(postgresCtx, a.settings.Postgres)
	if err != nil {
		return fmt.Errorf("can't connect to postgres: %v", err)
	}

	err = database.Migrate(postgres)
	if err != nil {
		return fmt.Errorf("can't migrate: %v", err)
	}

	a.userRepository = user.NewSqlRepository(postgres)
	return nil
}

func (a *App) InitServices() error {
	privateKey, err := fromBase64(a.settings.JWT.PrivateKey)
	if err != nil {
		return fmt.Errorf("can't parse private key: %v", err)
	}

	publicKey, err := fromBase64(a.settings.JWT.PublicKey)
	if err != nil {
		return fmt.Errorf("can't parse public key: %v", err)
	}

	ed25519, err := alg.NewEd25519(privateKey, publicKey)
	if err != nil {
		return fmt.Errorf("can't create Ed25519: %v", err)
	}

	jwt.Register(alg.EdDSA, ed25519, ed25519)
	a.encryptor = hash.NewBCrypt()
	a.userService = service.NewUserServiceImpl(a.userRepository, a.encryptor)
	return nil
}

func (a *App) Start(ctx service.ContextProvider) {
	a.server = api.NewServer(fmt.Sprintf(":%d", a.settings.Port), ctx, a.userService)
	a.server.Start()
}

func (a *App) Shutdown() error {
	ctx, cancelCtx := a.timeoutContext(time.Second * 15)
	defer cancelCtx()

	return a.server.Shutdown(ctx)
}

func fromBase64(data string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(data)
}
