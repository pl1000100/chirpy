package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pl1000100/chirpy/internal/database/auth"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type jsonRequestData struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	var requestData jsonRequestData

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, `{"error":"Could't decode"}`, http.StatusBadRequest)
		return
	}
	expires := 3600
	if 0 < requestData.ExpiresInSeconds && requestData.ExpiresInSeconds < 3600 {
		expires = requestData.ExpiresInSeconds
	}
	expiresDuration, err := time.ParseDuration(fmt.Sprintf("%ds", expires))
	if err != nil {
		http.Error(w, `{"error":"Could't parse expires"}`, http.StatusInternalServerError)
		return
	}
	user, err := cfg.db.GetUserByEmail(r.Context(), requestData.Email)
	if err != nil {
		http.Error(w, `{"error":"Could't get user from database"}`, http.StatusInternalServerError)
		return
	}

	if err := auth.CheckPasswordHash(user.HashedPassword, requestData.Password); err != nil {
		http.Error(w, `{"error":"Incorrect email or password"}`, http.StatusUnauthorized)
		return
	}
	generatedToken, err := auth.MakeJWT(user.ID, cfg.jwt_secret, expiresDuration)
	if err != nil {
		http.Error(w, `{"error":"Could't make JWT"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(
		User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token:     generatedToken,
		},
	); err != nil {
		http.Error(w, `{"error": "Could not encode use data"}`, http.StatusInternalServerError)
		return
	}

}
