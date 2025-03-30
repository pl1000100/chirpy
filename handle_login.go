package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pl1000100/chirpy/internal/auth"
	"github.com/pl1000100/chirpy/internal/database"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	const expireTimeToken = 1 * time.Hour
	const expireTimeRefreshToken = 60 * 24 * time.Hour

	// decode request
	var requestData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, `{"error":"Could't decode"}`, http.StatusBadRequest)
		return
	}

	// get user from database
	user, err := cfg.db.GetUserByEmail(r.Context(), requestData.Email)
	if err != nil {
		http.Error(w, `{"error":"Could't get user from database"}`, http.StatusInternalServerError)
		return
	}

	// check password
	if err := auth.CheckPasswordHash(user.HashedPassword, requestData.Password); err != nil {
		http.Error(w, `{"error":"Incorrect email or password"}`, http.StatusUnauthorized)
		return
	}

	// generate token
	generatedToken, err := auth.MakeJWT(user.ID, cfg.jwt_secret, expireTimeToken)
	if err != nil {
		http.Error(w, `{"error":"Could't make JWT"}`, http.StatusInternalServerError)
		return
	}

	// generate and store refresh token
	generatedRefreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		http.Error(w, `{"error":"Could't make refresh token"}`, http.StatusInternalServerError)
		return
	}
	databaseRefreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     generatedRefreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(expireTimeRefreshToken),
	})
	if err != nil {
		http.Error(w, `{"error":"Could't store refresh token"}`, http.StatusInternalServerError)
		return
	}

	// send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(
		User{
			ID:           user.ID,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Email:        user.Email,
			Token:        generatedToken,
			RefreshToken: databaseRefreshToken.Token,
			IsChirpyRed:  user.IsChirpyRed,
		},
	); err != nil {
		http.Error(w, `{"error": "Could not encode use data"}`, http.StatusInternalServerError)
		return
	}

}
