package handlers

import (
	"encoding/json"
	"errors"
	"literate-barnacle/database"
	"literate-barnacle/service/ctx"
	"net/http"

	"go.uber.org/zap"
)

func WriteJson(w http.ResponseWriter, log *zap.Logger, status int, response interface{}) {
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error("failed", zap.Error(err))
	}
}

func ProcessResponse(w http.ResponseWriter, log *zap.Logger, response interface{}, err error) {
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
	}
}
