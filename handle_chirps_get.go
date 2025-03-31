package main

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleChirpsGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get query parameters
	authorIdQuery := r.URL.Query().Get("author_id")
	sortQuery := r.URL.Query().Get("sort")

	allChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		http.Error(w, `{"error": "Could not get chirps"}`, http.StatusInternalServerError)
		return
	}

	jsonChirps := []Chirp{}
	for _, ch := range allChirps {
		if authorIdQuery == "" || authorIdQuery == ch.UserID.String() {
			jsonChirps = append(jsonChirps, Chirp{
				ID:        ch.ID,
				CreatedAt: ch.CreatedAt,
				UpdatedAt: ch.UpdatedAt,
				Body:      ch.Body,
				UserID:    ch.UserID,
			})
		}
	}

	if sortQuery == "desc" {
		sort.Slice(jsonChirps, func(i, j int) bool {
			return jsonChirps[i].CreatedAt.After(jsonChirps[j].CreatedAt) // desc
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
