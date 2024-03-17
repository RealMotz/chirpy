package main

import (
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

func readinessHandler(w http.ResponseWriter, r *http.Request) { 
  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("OK"))
}

func main() {
  port := ":8080" 
  directoryPath := "."
  mux := http.NewServeMux()
  crsMux := middlewareCors(mux)

  server := &http.Server{
    Addr: port,
    Handler: crsMux,
  }

  mux.Handle("/app/*", http.StripPrefix("/app/", http.FileServer(http.Dir(directoryPath))))
  mux.HandleFunc("/healthz", readinessHandler)
  log.Printf("Serving on port %s", port)

  err := server.ListenAndServe()
  if err != nil {
    log.Fatal(err)
  }
}
