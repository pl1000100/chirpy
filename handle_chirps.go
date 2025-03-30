package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pl1000100/chirpy/internal/auth"
	"github.com/pl1000100/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
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

func (cfg *apiConfig) handleChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		Body string `json:"body"`
	}

	w.Header().Set("Content-Type", "application/json")

	var jsonData returnVals
	if err := json.NewDecoder(r.Body).Decode(&jsonData); err != nil {
		http.Error(w, `{"error": "Couldn't decode"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {

		http.Error(w, `{"error": "Couldn't authorize"}`, http.StatusBadRequest)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwt_secret)
	if err != nil {
		http.Error(w, `{"error": "Unauthorize"}`, http.StatusUnauthorized)
		return
	}

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
		UserID: userID,
	}
	createdChirp, err := cfg.db.CreateChirp(r.Context(), createChirpParams)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Chirp could not be created, %v"}`, err), http.StatusInternalServerError)
		http.Error(w, `{"error": "Chirp could not be created"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	respData := Chirp{
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

func (cfg *apiConfig) handleChirpsGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		http.Error(w, `{"error": "Could not get chirps"}`, http.StatusInternalServerError)
		return
	}
	jsonChirps := []Chirp{}
	for _, ch := range allChirps {
		jsonChirps = append(jsonChirps, Chirp{
			ID:        ch.ID,
			CreatedAt: ch.CreatedAt,
			UpdatedAt: ch.UpdatedAt,
			Body:      ch.Body,
			UserID:    ch.UserID,
		})
	}
	if err := json.NewEncoder(w).Encode(jsonChirps); err != nil {
		http.Error(w, `{"error": "Could not encode chirp"}`, http.StatusInternalServerError)
		return
	}
}

func (cfg *apiConfig) handleChirpsGetOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		http.Error(w, `{"error":"Couldn't parse id"}`, http.StatusBadRequest)
		return
	}
	dbChirp, err := cfg.db.GetOneChirp(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error":"Couldn't get chirp"}`, http.StatusNotFound)
		return
	}
	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	if err := json.NewEncoder(w).Encode(chirp); err != nil {
		http.Error(w, `{"error": "Could not encode chirp"}`, http.StatusInternalServerError)
		return
	}
}

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
