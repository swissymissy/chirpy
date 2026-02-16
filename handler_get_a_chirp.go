package main 

import (
	"net/http"
	"fmt"
	"errors"
	"database/sql"

	"github.com/google/uuid"
)


func (apicfg *apiConfig) handlerGetAChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")		// extract ID string from the path
	chirpID, err := uuid.Parse(chirpIDStr)		// convert the ID type string into UUID type
	if err != nil {
		respondWithError(w , http.StatusBadRequest, "invalid chirp ID")
		return
	}
	
	// retrieve a single row
	oneChirp, err := apicfg.DB.GetAChirp(r.Context(), chirpID)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Printf("Error getting a row from chirps: %s\n", err)
		msg := "Chirp does not exist"
		respondWithError(w, 404, msg)
		return
	} else if err != nil {
		fmt.Printf("Error getting row from chirps: %s\n", err)
		msg := "Can't get chirp"
		respondWithError(w, 500, msg)
		return
	}

	// convert to right response format
	resFmt := resFormat{
		ID: oneChirp.ID,
		CreatedAt: oneChirp.CreatedAt,
		UpdatedAt: oneChirp.UpdatedAt,
		Body: oneChirp.Body,
		UserID: oneChirp.UserID,
	}
	respondWithJSON(w, 200, resFmt)
	return
}