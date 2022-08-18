package handlers

import (
	"encoding/json"
	"literate-barnacle/service"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func LoginHandler(rawLog *zap.Logger, user service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rayId := uuid.NewString()
		log := rawLog.With(zap.String("rayId", rayId), zap.String("method", "auth/login"))

		request := service.LoginRequest{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Warn("bad request", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		response, err := user.Login(r.Context(), request)
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
