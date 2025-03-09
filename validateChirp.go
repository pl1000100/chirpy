package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		Body string `json:"body"`
	}
	w.Header().Set("Content-Type", "application/json")
	var jsonData returnVals
	if err := json.NewDecoder(r.Body).Decode(&jsonData); err != nil {
		http.Error(w, `{"error": "Something went wrong"}`, http.StatusBadRequest)
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

	newBody := strings.Split(jsonData.Body, " ")

	for _, wo := range words {
		for j, word := range newBody {
			if strings.ToLower(word) == wo {
				newBody[j] = "****"
			}
		}

	}

	w.WriteHeader(http.StatusOK)
	// w.Write([]byte(`{"valid": true}`))
	w.Write([]byte(fmt.Sprintf(`{"cleaned_body": "%s"}`, strings.Join(newBody, " "))))
}
