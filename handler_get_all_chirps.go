package main 

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (apicfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	var responseFmt []resFormat
	
	// grab query parameter from URL, there is one
	authorId := r.URL.Query().Get("author_id")

	if authorId == "" {
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
	} else {
		// convert ID string to UUID
		userID, err := uuid.Parse(authorId)
		if err != nil {
			respondWithError(w, 404, "Invalid ID")
			return
		}
		
		// get list of chirps created by this userID
		chirpList, err := apicfg.DB.GetAllChirpsFromUserID(r.Context(), userID)
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
	}
}