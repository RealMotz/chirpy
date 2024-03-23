package main

import (
	"fmt"
	"log"
	"net/http"
)

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Add("Content-Type", "text/plain; charset=utf-8")
  w.WriteHeader(http.StatusOK)
  w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
  cfg.fileserverHits = 0
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("Hits reset to 0"))
}

func main() {
	port := ":8080"
	directoryPath := "."
	mux := http.NewServeMux()
	crsMux := middlewareCors(mux)
	config := apiConfig{fileserverHits: 0}

	server := &http.Server{
		Addr:    port,
		Handler: crsMux,
	}

	mux.Handle("/app/*", config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(directoryPath)))))
	mux.HandleFunc("/healthz", readinessHandler)
	mux.HandleFunc("/metrics", config.metricsHandler)
	mux.HandleFunc("/reset", config.resetHandler)
	log.Printf("Serving on port %s", port)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
