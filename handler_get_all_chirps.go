package main 

import (
	"fmt"
	"net/http"


)

func (apicfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	var responseFmt []resFormat
	
	chirpList, err := apicfg.DB.GetAllChirps(r.Context())
	if err != nil {
		fmt.Printf("Error getting all chirps from db: %s\n", err)
		msg := "Couldn't get all chirps!"
		respondWithError(w, 500, msg)
		return 
	}

	// writing each item from chirpList to response format
	for _, eachChirp := range chirpList{
		responseFmt = append(responseFmt, resFormat{
				ID: eachChirp.ID,
				CreatedAt: eachChirp.CreatedAt,
				UpdatedAt: eachChirp.UpdatedAt,
				Body: eachChirp.Body,
				UserID: eachChirp.UserID,
			})
	}

	respondWithJSON(w, 200, responseFmt)
	return 
}