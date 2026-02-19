package main

import (
	"net/http"
	"encoding/json"
	
	"github.com/google/uuid"
	"github.com/swissymissy/chirpy/internal/auth"
)

// function used for communitcating with third-party server: Polka
func (apicfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {
	// validate the API key
	apiKey, err := auth.GetAPIKey(r.Header) 
	if err != nil {
		w.WriteHeader(401)
		return
	}
	
	if apiKey != apicfg.polkaKey {
		w.WriteHeader(401)
		return
	}

	type WebhookEvent struct {
		Event string 	`json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		}`json:"data"`
	}

	//decode the request body
	decoder := json.NewDecoder(r.Body)
	var data WebhookEvent
	err = decoder.Decode(&data)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// check the event
	if data.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	// parse the ID string into UUID
	userID, err := uuid.Parse(data.Data.UserID)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	// event is user.upgrade -> update in db
	_, err = apicfg.DB.UpgradeUserToChirpyRed(r.Context(), userID)
	if err != nil {
		w.WriteHeader(404)
		return
	} 
	w.WriteHeader(204)

}