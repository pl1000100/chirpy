package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/pl1000100/chirpy/internal/auth"
	"github.com/pl1000100/chirpy/internal/database"
)

func (cfg *apiConfig) handlePolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type polkaUserUpgraded struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	w.Header().Set("Content-Type", "application/json")

	// get request data
	reqData := polkaUserUpgraded{}
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, `{"error":"Couldn't decode"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// auth
	reqApiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || reqApiKey != cfg.polka_key {
		http.Error(w, `{"error":"Couldn't get api key or bad key"}`, http.StatusUnauthorized)
		return
	}

	// handle events
	switch reqData.Event {

	case "user.upgraded":
		userID, err := uuid.Parse(reqData.Data.UserID)
		if err != nil {
			http.Error(w, `{"error":"Couldn't parse user"}`, http.StatusInternalServerError)
			return
		}
		if err := cfg.db.UpdateUserChirpyRedByID(
			r.Context(),
			database.UpdateUserChirpyRedByIDParams{
				ID:          userID,
				IsChirpyRed: true,
			},
		); err != nil {
			http.Error(w, `{"error":"Couldn't find user"}`, http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		w.WriteHeader(http.StatusNoContent)

	}

}
