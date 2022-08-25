package handlers

import (
	"encoding/json"
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
