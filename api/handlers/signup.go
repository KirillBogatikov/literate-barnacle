package handlers

import (
	"encoding/json"
	"literate-barnacle/service"
	"log"
	"net/http"
)

func SignUpHandler(user service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := service.SignUpRequest{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		response, err := user.Signup(r.Context(), request)
		if err != nil {
			log.Printf("signup failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("signup failed: %v", err)
		}
	}
}
