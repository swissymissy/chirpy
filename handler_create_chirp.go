package main 

import (
	"fmt"
	"net/http"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/swissymissy/chirpy/internal/database"
	"github.com/swissymissy/chirpy/internal/auth"
)

type chirpMsg struct {
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}	

// response format
type resFormat struct {
	ID        uuid.UUID		`json:"id"`
	CreatedAt time.Time		`json:"created_at"`
	UpdatedAt time.Time		`json:"updated_at"`
	Body      string		`json:"body"`
	UserID    uuid.UUID		`json:"user_id"`
}

func (apicfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	
	// decode body req into json bytes
	decoder := json.NewDecoder(r.Body)
	var chrpmsg chirpMsg
	err := decoder.Decode(&chrpmsg)	// write to params after decoding
	if err != nil {
		fmt.Printf("Error decoding request body: %s", err)
		msg := "Something went wrong"
		respondWithError(w, 500, msg )
		return
	}	
	
	// get the token from request body
	token, err := auth.GetBearerToken(r.Header) 
	if err != nil {
		fmt.Printf("Error getting token from header: %s", err)
		msg := "Invalid token"
		respondWithError(w, 401, msg)
		return
	}

	// validate the token
	userId, err := auth.ValidateJWT(token, apicfg.jwt_secret)
	if err != nil {
		fmt.Printf("Invalid token: %w\n", err)
		msg := "Invalid token"
		respondWithError(w, 401, msg)
		return
	}

	// validate the repsonse body
	err = ValidateChirp(&chrpmsg)
	if err != nil {
		fmt.Printf("Invalid request body: %s\n", err)
		msg := "Invalid chirp"
		respondWithError(w, 400, msg)
		return 
	}

	// create ChirpParams
	cp := database.CreateChirpParams{
		Body: chrpmsg.Body,
		UserID: userId,
	}
	// create new chirp msg in db
	newChirp, err := apicfg.DB.CreateChirp(r.Context(), cp)
	if err != nil {
		fmt.Printf("Error adding new chirp to db: %s\n", err)
		msg := "Can't create chirp"
		respondWithError(w, 500, msg)
		return 
	}

	rf := resFormat{
		ID: newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body: newChirp.Body,
		UserID: newChirp.UserID,
	}
	respondWithJSON(w, 201, rf)
	return
}