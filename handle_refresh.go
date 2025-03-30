package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pl1000100/chirpy/internal/auth"
)

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	const expireTimeNewToken = 1 * time.Hour
	w.Header().Set("Content-Type", "application/json")

	// get refresh token
	reqRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, `{"error": "Could not get token from header"}`, http.StatusUnauthorized)
		return
	}

	// get and validate refresh token from database
	dbRefreshToken, err := cfg.db.GetRefreshToken(r.Context(), reqRefreshToken)
	if err != nil {
		http.Error(w, `{"error": "Could not get token from database"}`, http.StatusUnauthorized)
		return
	}
	if dbRefreshToken.ExpiresAt.Before(time.Now()) || dbRefreshToken.RevokedAt.Valid {
		http.Error(w, `{"error": "Token expired or revoked"}`, http.StatusUnauthorized)
		return
	}

	// get user
	dbUser, err := cfg.db.GetUserByRereshToken(r.Context(), dbRefreshToken.Token)
	if err != nil {
		http.Error(w, `{"error": "Could not get user from database"}`, http.StatusInternalServerError)
		return
	}

	// create new token
	newToken, err := auth.MakeJWT(dbUser.ID, cfg.jwt_secret, expireTimeNewToken)
	if err != nil {
		http.Error(w, `{"error": "Could generate new token"}`, http.StatusInternalServerError)
		return
	}

	// respond with new token
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(struct {
		Token string `json:"token"`
	}{
		Token: newToken,
	},
	); err != nil {
		http.Error(w, `{"error": "Couldn't encode token"}`, http.StatusInternalServerError)
		return
	}

}

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {
	// get refresh token
	reqRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, `{"error": "Could not get token from header"}`, http.StatusUnauthorized)
		return
	}

	// get and validate refresh token from database
	dbRefreshToken, err := cfg.db.GetRefreshToken(r.Context(), reqRefreshToken)
	if err != nil {
		http.Error(w, `{"error": "Could not get token from database"}`, http.StatusUnauthorized)
		return
	}
	if dbRefreshToken.ExpiresAt.Before(time.Now()) || dbRefreshToken.RevokedAt.Valid {
		http.Error(w, `{"error": "Token expired or revoked"}`, http.StatusUnauthorized)
		return
	}

	// Revoke token
	if err := cfg.db.RevokeRefreshToken(r.Context(), dbRefreshToken.Token); err != nil {
		http.Error(w, `{"error": "Could not revoke token"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
