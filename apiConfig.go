package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Raghart/go-servers/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newValue := cfg.fileserverHits.Add(1)
		cfg.fileserverHits.Store(newValue)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	totalHitsStr := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())

	w.Write([]byte(totalHitsStr))
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)

	if platform := os.Getenv("PLATFORM"); platform != "dev" {
		jsonResponseError(w, 403, "Forbidden")
		return
	}

	err := cfg.db.ResetUsers(context.Background())

	if err != nil {
		jsonResponseError(w, 400, fmt.Sprintf("there was an error while trying to reset the database: %v", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
