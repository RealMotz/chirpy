package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/RealMotz/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	db             database.DataBase
}

func main() {
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	port := ":8080"
	directoryPath := "."
	mux := http.NewServeMux()
	crsMux := middlewareCors(mux)

	var mtx sync.RWMutex
	config := apiConfig{
		fileserverHits: 0,
		db: database.DataBase{
			Name: "database.json",
			Mux:  &mtx,
		},
	}

	if *debug {
		os.Remove(config.db.Name)
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
	mux.HandleFunc("/api/reset", config.resetHandler)
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", config.metricsHandler)
	mux.HandleFunc("GET /api/chirps", config.getChirps)
	mux.HandleFunc("GET /api/chirps/{id}", config.getChirp)
	mux.HandleFunc("POST /api/chirps", config.createChirp)
	mux.HandleFunc("POST /api/users", config.createUser)
	log.Printf("Serving on port %s", port)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
