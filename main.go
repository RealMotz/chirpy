package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/RealMotz/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	db             database.DataBase
}

type errorResponse struct {
	Error string `json:"error"`
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	body := `
      <html>
        <body>
            <h1>Welcome, Chirpy Admin</h1>
            <p>Chirpy has been visited %d times!</p>
        </body>
      </html>
  `
	w.Write([]byte(fmt.Sprintf(body, cfg.fileserverHits)))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decodedChirp := database.Chirp{}
	err := decoder.Decode(&decodedChirp)
	if err != nil {
		handleJsonResponse(w, http.StatusInternalServerError, errorResponse{
			Error: "Cannot decode parameters",
		})
		return
	}

	err = validateChirp(decodedChirp)
	if err != nil {
		handleJsonResponse(w, http.StatusBadRequest, errorResponse{
			Error: err.Error(),
		})
		return
	}

	err = cfg.db.CreateIfNotExits()
	if err != nil {
		handleJsonResponse(w, http.StatusInternalServerError, errorResponse{
			Error: err.Error(),
		})
		return
	}

	chirp, err := cfg.db.CreateChirp(cleanBody(decodedChirp.Body))
	if err != nil {
		handleJsonResponse(w, http.StatusInternalServerError, errorResponse{
			Error: err.Error(),
		})
		return
	}
	handleJsonResponse(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	err := cfg.db.CreateIfNotExits()
	if err != nil {
		handleJsonResponse(w, http.StatusInternalServerError, errorResponse{
			Error: err.Error(),
		})
		return
	}

	chirps, err := cfg.db.GetChirps()
	if err != nil {
		handleJsonResponse(w, http.StatusInternalServerError, errorResponse{
			Error: err.Error(),
		})
		return
	}

	handleJsonResponse(w, http.StatusOK, chirps)
}

func validateChirp(chirp database.Chirp) error {
	if len(chirp.Body) > 140 {
		return errors.New("chirp is too long")
	}
	return nil
}

func handleJsonResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error enconding json %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}
	w.WriteHeader(statusCode)
	w.Write(data)
}

func cleanBody(body string) string {
	profanities := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	words := strings.Split(body, " ")
	for i := 0; i < len(words); i++ {
		for _, profanity := range profanities {
			if profanity == strings.ToLower(words[i]) {
				words[i] = "****"
				break
			}
		}
	}
	return strings.Join(words, " ")
}

func main() {
	port := ":8080"
	directoryPath := "."
	mux := http.NewServeMux()
	crsMux := middlewareCors(mux)

	var mtx sync.RWMutex
	config := apiConfig{
		fileserverHits: 0,
		db: database.DataBase{
			Name: "chirps.json",
			Mux:  &mtx,
		},
	}

	server := &http.Server{
		Addr:    port,
		Handler: crsMux,
	}

	mux.Handle(
		"/app/*",
		config.middlewareMetricsInc(
			http.StripPrefix("/app", http.FileServer(http.Dir(directoryPath))),
		),
	)
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", config.metricsHandler)
	mux.HandleFunc("/api/reset", config.resetHandler)
	mux.HandleFunc("GET /api/chirps", config.getChirps)
	mux.HandleFunc("POST /api/chirps", config.createChirp)
	log.Printf("Serving on port %s", port)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
