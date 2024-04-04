package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handleErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	handleJsonResponse(w, statusCode, errorResponse{
		Error: err.Error(),
	})
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
