package main

import (
	"encoding/json"
	"net/http"

	"github.com/pl1000100/chirpy/internal/auth"
	"github.com/pl1000100/chirpy/internal/database"
)

func (cfg *apiConfig) handleUsersUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// get request data
	var reqData struct {
		Email   string `json:"email"`
		Pasword string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, `{"error": "Can't decode request's body"}`, http.StatusBadRequest)
		return
	}

	// authenticate user
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, `{"error": "Can't get token"}`, http.StatusUnauthorized)
		return
	}
	userID, err := auth.ValidateJWT(accessToken, cfg.jwt_secret)
	if err != nil {
		http.Error(w, `{"error": "Can't validate JWT"}`, http.StatusUnauthorized)
		return
	}

	// hash password
	hashedPassword, err := auth.HashPassword(reqData.Pasword)
	if err != nil {
		http.Error(w, `{"error": "Can't hash password"}`, http.StatusBadRequest)
		return
	}

	// update users data
	updatedUser, err := cfg.db.UpdateUserByToken(
		r.Context(),
		database.UpdateUserByTokenParams{
			ID:             userID,
			Email:          reqData.Email,
			HashedPassword: hashedPassword,
		},
	)
	if err != nil {
		http.Error(w, `{"error": "Can't update database"}`, http.StatusInternalServerError)
		return
	}

	// respond
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(User{
		ID:          updatedUser.ID,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
		Email:       updatedUser.Email,
		IsChirpyRed: updatedUser.IsChirpyRed,
	}); err != nil {
		http.Error(w, `{"error": "Can't encode data"}`, http.StatusInternalServerError)
		return
	}
}
