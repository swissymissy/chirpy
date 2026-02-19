package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	
	"github.com/google/uuid"
)


func (apicfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {
	type WebhookEvent struct {
		Event string 	`json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		}`json:"data"`
	}

	//decode the request body
	decoder := json.NewDecoder(r.Body)
	var data WebhookEvent
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Printf("Error decoding body request: %s", err)
		msg := "Something went wrong"
		respondWithError(w, 500, msg )
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