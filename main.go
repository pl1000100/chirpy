package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pl1000100/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	jwt_secret     string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	const port = "8080"
	const filePathRoot = "."
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	jwt_secret := os.Getenv("JWT_SECRET")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	dbQueries := database.New(db)
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		jwt_secret:     jwt_secret,
	}

	mux := http.NewServeMux()

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	
	
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerMetricsReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	
	mux.HandleFunc("POST /api/users", apiCfg.handleUsersCreate)
	mux.HandleFunc("PUT /api/users", apiCfg.handleUsersUpdate)
	
	mux.HandleFunc("POST /api/chirps", apiCfg.handleChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handleChirpsGetAll)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handleChirpsGetOne)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handleChirpsDeleteOne)
	
	mux.HandleFunc("POST /api/login", apiCfg.handleLogin)
	
	mux.HandleFunc("POST /api/refresh", apiCfg.handleRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handleRevoke)
	
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlePolkaWebhooks)
	
	svr := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	svr.ListenAndServe()
}
