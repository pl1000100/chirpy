package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pl1000100/chirpy/internal/auth"
	"github.com/pl1000100/chirpy/internal/database"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handleUsersCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type returnVals struct {
		Email   string `json:"email"`
		Pasword string `json:"password"`
	}

	var params returnVals
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, `{"error": "Can't decode request's body"}`, http.StatusBadRequest)
		return
	}

	hashed, err := auth.HashPassword(params.Pasword)
	if err != nil {
		http.Error(w, `{"error": "Can't hash password"}`, http.StatusBadRequest)
		return
	}

	user, err := cfg.db.CreateUser(
		r.Context(),
		database.CreateUserParams{
			Email:          params.Email,
			HashedPassword: hashed,
		},
	)

	if err != nil {
		http.Error(w, `{"error": "Can't create user"}`, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)

	res, err := json.Marshal(
		User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
	)
	if err != nil {
		http.Error(w, `{"error": "Can't marshal data"}`, http.StatusBadRequest)
		return
	}
	w.Write(res)
}

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
