package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	port := "8080"
	mux := http.NewServeMux()
	apiCfg := &apiConfig{}

	mux.HandleFunc("/healthz", healthHandler)
	mux.Handle("/metrics", apiCfg)
	mux.HandleFunc("/reset", apiCfg.resetMetrics)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newValue := cfg.fileserverHits.Add(1)
		cfg.fileserverHits.Store(newValue)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	totalHits := cfg.fileserverHits.Load()

	w.WriteHeader(http.StatusOK)
	totalHitsStr := fmt.Sprintf("Hits: %d", totalHits)
	w.Write([]byte(totalHitsStr))
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
