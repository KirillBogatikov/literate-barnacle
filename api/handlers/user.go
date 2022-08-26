package handlers

import (
	"encoding/json"
	"literate-barnacle/service/ctx"
	"literate-barnacle/service/models"
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
		defer log.Debug("request finished", zap.String("executionTime", time.Since(start).String()))

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
		ProcessResponse(w, log, response, err)
	}
}

func UpdateUserHandler(rawLog *zap.Logger, service user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log := rawLog.With(zap.String("method", "user/update"))
		defer log.Debug("request finished", zap.String("executionTime", time.Since(start).String()))

		c, err := ctx.GetContext(r, true)
		if err != nil {
			log.Warn("forbidden", zap.Error(err))
			w.WriteHeader(http.StatusForbidden)
			return
		}

		var domainUser models.User
		if err = json.NewDecoder(r.Body).Decode(&domainUser); err != nil {
			log.Warn("bad request", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		response, err := service.Update(c, domainUser)
		ProcessResponse(w, log, response, err)
	}
}

func DeleteUserHandler(rawLog *zap.Logger, service user.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log := rawLog.With(zap.String("method", "user/delete"))
		defer log.Debug("request finished", zap.String("executionTime", time.Since(start).String()))

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

		response, err := service.Delete(c, userId)
		ProcessResponse(w, log, response, err)
	}
}
