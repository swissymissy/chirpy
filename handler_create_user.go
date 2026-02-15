package main

import (
	"net/http"
	"time"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)


type User struct {
	ID 			uuid.UUID 	`json:"id"`
	CreatedAt	time.Time 	`json:"created_at"`
	UpdatedAt	time.Time 	`json:"updated_at"`
	Email		string 		`json:"email"`
}

func (apicfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	
	type reqEmail struct {
		Email	string 	`json:"email"`
	}

	// decode body req into json bytes
	decoder := json.NewDecoder(r.Body)
	var rE reqEmail
	err := decoder.Decode(&rE)
	if err != nil {
		fmt.Printf("Error decoding body request: %s", err)
		msg := "Something went wrong"
		respondWithError(w, 500, msg )
		return
	}

	// inserting user to db
	newUsr, err := apicfg.DB.CreateUser(r.Context(), rE.Email)
	if err != nil {
		fmt.Printf("Error creating new user in db: %s", err)
		msg := "Something went wrong"
		respondWithError(w, 500, msg )
		return
	}

	// encode user struct into json
	resUsr := User{
		ID: newUsr.ID,
		CreatedAt: newUsr.CreatedAt,
		UpdatedAt: newUsr.UpdatedAt,
		Email: newUsr.Email,
	}
	respondWithJSON(w, 201, resUsr)
	return
} 