package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pl1000100/chirpy/internal/database"
)

func (cfg *apiConfig) handleChirps(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	w.Header().Set("Content-Type", "application/json")

	var jsonData returnVals
	if err := json.NewDecoder(r.Body).Decode(&jsonData); err != nil {
		http.Error(w, `{"error": "Couldn't decode"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if len(jsonData.Body) > 140 {
		http.Error(w, `{"error": "Chirp is too long"}`, http.StatusBadRequest)
		return
	}

	words := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	newBody := filterWords(jsonData.Body, words)

	createChirpParams := database.CreateChirpParams{
		Body:   newBody,
		UserID: jsonData.UserID,
	}
	createdChirp, err := cfg.db.CreateChirp(r.Context(), createChirpParams)
	if err != nil {
		http.Error(w, `{"error": "Chirp could not be created"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	respData := struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}{
		ID:        createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt,
		UpdatedAt: createdChirp.UpdatedAt,
		Body:      createdChirp.Body,
		UserID:    createdChirp.UserID,
	}
	if err := json.NewEncoder(w).Encode(respData); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func filterWords(body string, words []string) string {
	newBodySplitted := strings.Split(body, " ")
	for _, wo := range words {
		for j, word := range newBodySplitted {
			if strings.ToLower(word) == wo {
				newBodySplitted[j] = "****"
			}
		}
	}
	return strings.Join(newBodySplitted, " ")
}
