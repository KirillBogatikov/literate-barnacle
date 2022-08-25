package handlers

import (
	"encoding/json"
	"literate-barnacle/service/ctx"
	"literate-barnacle/service/user"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func LoginHandler(rawLog *zap.Logger, service user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log := rawLog.With(zap.String("method", "auth/login"))

		request := user.LoginRequest{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Warn("bad request", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		c, _ := ctx.GetContext(r, false)
		response, err := service.Login(c, request)
		if err != nil {
			log.Error("login failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if response.IsSuccess() {
			log.Info("user authorized", zap.String("login", request.Login))
		} else {
			log.Info("request validation failed")
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Error("login failed", zap.Error(err))
		}

		log.Debug("request finished", zap.String("executionTime", time.Since(start).String()))
	}
}
