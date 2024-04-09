package main

import (
	"encoding/json"
	"net/http"

	"github.com/RealMotz/chirpy/internal/auth"
	"github.com/RealMotz/chirpy/internal/database"
)

type webhookRequest struct {
	Event string         `json:"event"`
	Data  map[string]int `json:"data"`
}

func (cfg *apiConfig) webhooks(w http.ResponseWriter, r *http.Request) {
	token, err := auth.FetchAuthHeader(r.Header.Get("Authorization"), "ApiKey")
	if err != nil {
		handleErrorResponse(w, http.StatusUnauthorized, err)
		return
	}

	if token != string(cfg.polkaApiKey) {
		handleErrorResponse(w, http.StatusUnauthorized, err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	webhookReq := webhookRequest{}
	err = decoder.Decode(&webhookReq)
	if err != nil {
		handleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	if webhookReq.Event != database.UpgradedEvent {
		handleJsonResponse(w, http.StatusOK, struct{}{})
		return
	}

	for key := range webhookReq.Data {
		err = cfg.db.AddMembership(webhookReq.Data[key])
		if err != nil {
			handleErrorResponse(w, http.StatusNotFound, err)
			return
		}
	}

	handleJsonResponse(w, http.StatusOK, struct{}{})

}
