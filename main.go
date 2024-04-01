package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

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

func (cfg *apiConfig) validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	type validResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		handleJsonResponse(w, http.StatusInternalServerError, errorResponse{
			Error: "Cannot decode parameters",
		})
		return
	}

	if len(params.Body) > 140 {
		handleJsonResponse(w, http.StatusBadRequest, errorResponse{
			Error: "Chirp is too long",
		})
		return
	}

	handleJsonResponse(w, http.StatusOK, validResponse{
		CleanedBody: cleanBody(params.Body),
	})
}

func handleJsonResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error enconding json %s", err)
		w.WriteHeader(500)
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
	config := apiConfig{fileserverHits: 0}

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
	mux.HandleFunc("/api/validate_chirp", config.validateChirp)
	log.Printf("Serving on port %s", port)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
