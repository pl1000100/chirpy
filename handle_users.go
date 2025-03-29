package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pl1000100/chirpy/internal/database"
	"github.com/pl1000100/chirpy/internal/database/auth"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handleUsersCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type returnVals struct {
		Email   string `json:"email"`
		Pasword string `json:"password"`
	}

	var params returnVals
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, `{"error": "Can't decode request's body"}`, http.StatusBadRequest)
		return
	}

	hashed, err := auth.HashPassword(params.Pasword)
	if err != nil {
		http.Error(w, `{"error": "Can't hash password"}`, http.StatusBadRequest)
		return
	}

	user, err := cfg.db.CreateUser(
		r.Context(),
		database.CreateUserParams{
			Email:          params.Email,
			HashedPassword: hashed,
		},
	)

	if err != nil {
		http.Error(w, `{"error": "Can't create user"}`, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)

	res, err := json.Marshal(
		User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	)
	if err != nil {
		http.Error(w, `{"error": "Can't marshal data"}`, http.StatusBadRequest)
		return
	}
	w.Write(res)
}
