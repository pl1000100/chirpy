package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleUsers(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		Email string `json:"email"`
	}

	var params returnVals
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, `{"error": "Can't decode request's body"}`, http.StatusBadRequest)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		http.Error(w, `{"error": "Can't create user"}`, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	response := struct {
		ID        uuid.UUID `json:"id"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, `{"error": "Can't marshal data"}`, http.StatusBadRequest)
		return
	}
	w.Write(res)
}
