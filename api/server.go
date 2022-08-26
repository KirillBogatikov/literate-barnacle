package api

import (
	"context"
	"errors"
	"literate-barnacle/api/handlers"
	"literate-barnacle/service/ctx"
	"literate-barnacle/service/user"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	server *http.Server
	log    *zap.Logger
}

func NewServer(
	address string,
	ctx ctx.ContextProvider,
	log *zap.Logger,
	user user.Service,
) Server {
	router := mux.NewRouter()

	router.Handle("/auth/login", handlers.LoginHandler(log, user)).Methods(http.MethodPost)
	router.Handle("/auth/signup", handlers.SignUpHandler(log, user)).Methods(http.MethodPost)

	router.Handle("/user/{userId}", handlers.GetUserHandler(log, user)).Methods(http.MethodGet)
	router.Handle("/user", handlers.GetUserHandler(log, user)).Methods(http.MethodGet)
	router.Handle("/user", handlers.UpdateUserHandler(log, user)).Methods(http.MethodPut)

	server := &http.Server{
		Addr:              address,
		Handler:           router,
		ReadHeaderTimeout: time.Second * 10,
		ReadTimeout:       time.Second * 15,
		WriteTimeout:      time.Second * 30,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx()
		},
	}

	return Server{
		server: server,
		log:    log,
	}
}

func (s *Server) Start() {
	go func() {
		err := s.server.ListenAndServe()
		switch {
		case err == nil:
			return
		case errors.Is(err, http.ErrServerClosed):
			return
		default:
			s.log.Error("server down", zap.Error(err))
		}
	}()
	s.log.Info("server up", zap.String("address", s.server.Addr))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
