package handlers

import (
	"errors"
	"literate-barnacle/database"
	"literate-barnacle/service/ctx"
	"literate-barnacle/service/user"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func GetUserHandler(rawLog *zap.Logger, service user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log := rawLog.With(zap.String("method", "user/get"))

		c, err := ctx.GetContext(r, true)
		if err != nil {
			log.Warn("forbidden", zap.Error(err))
			w.WriteHeader(http.StatusForbidden)
			return
		}

		vars := mux.Vars(r)
		id, ok := vars["userId"]

		var userId uuid.UUID
		if ok {
			userId, err = uuid.Parse(id)
			if err != nil {
				log.Warn("bad param", zap.Error(err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			userId = c.Authorization.UserId
		}

		response, err := service.Get(c, userId)
		switch {
		case err == nil:
			WriteJson(w, log, http.StatusOK, response)

		case errors.Is(err, database.ErrNotFound):
			log.Error("failed", zap.Error(err))
			WriteJson(w, log, http.StatusNotFound, response)

		case errors.Is(err, ctx.ErrUnauthorized):
			log.Error("failed", zap.Error(err))
			WriteJson(w, log, http.StatusUnauthorized, response)

		case errors.Is(err, ctx.ErrForbidden):
			log.Error("failed", zap.Error(err))
			WriteJson(w, log, http.StatusForbidden, response)

		default:
			log.Error("failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Debug("request finished", zap.String("executionTime", time.Since(start).String()))
	}
}
