package main 

import (
	"fmt"
	"net/http"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/swissymissy/chirpy/internal/auth"
)


func (apicfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	// check for access token header
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Printf("Error extracting token from header: %s\n", err)
		respondWithError(w , 401 , "Invalid Token")
		return
	}

	// validate the token
	userID, err := auth.ValidateJWT(accessToken, apicfg.jwt_secret)
	if err != nil {
		fmt.Printf("Error validating token: %s\n", err)
		respondWithError(w, 401, "Invalid token")
		return
	}

	chirpIDStr := r.PathValue("chirpID")		// extract ID string from URL
	chirpID, err := uuid.Parse(chirpIDStr)		// convert string to UUID type
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	// retrive chirp from table
	chirpInfo, err := apicfg.DB.GetAChirp(r.Context(), chirpID)
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

	// check if this chirp belong to the user
	if userID != chirpInfo.UserID {
		respondWithError(w , 403 , "Unauthorized")
		return
	}

	// delete chirp after checking authorization
	err = apicfg.DB.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 500, "Unable to delete chirp")
		return
	}
	w.WriteHeader(204)
}