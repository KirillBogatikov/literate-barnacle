package api

import (
	"context"
	"errors"
	"literate-barnacle/api/handlers"
	"literate-barnacle/service"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	server *http.Server
}

func NewServer(address string, ctx service.ContextProvider, user service.UserService) Server {
	router := mux.NewRouter()

	router.Handle("/auth/login", handlers.LoginHandler(user)).Methods(http.MethodPost)
	router.Handle("/auth/signup", handlers.SignUpHandler(user)).Methods(http.MethodPost)

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
			// TODO: log
		}
	}()
	log.Printf("server listening %s", s.server.Addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
