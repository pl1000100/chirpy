package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/pl1000100/chirpy/internal/auth"
)

func (cfg *apiConfig) handleChirpsDeleteOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// authenticate
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

	// get chirp
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		http.Error(w, `{"error":"Couldn't parse id"}`, http.StatusBadRequest)
		return
	}
	dbChirp, err := cfg.db.GetOneChirp(r.Context(), chirpID)
	if err != nil {
		http.Error(w, `{"error":"Couldn't get chirp from db"}`, http.StatusInternalServerError)
		return
	}

	//authorize and delete
	if dbChirp.UserID != userID {
		http.Error(w, `{"error":"Couldn't delete someone else's chirp"}`, http.StatusForbidden)
		return
	}
	if err := cfg.db.DeleteOneChirp(r.Context(), dbChirp.ID); err != nil {
		http.Error(w, `{"error":"Couldn't delete from database, not found"}`, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
