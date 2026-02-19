package main 

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/swissymissy/chirpy/internal/database"
)

func (apicfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	var responseFmt []resFormat
	
	// grab query parameter from URL, there is one
	authorId := r.URL.Query().Get("author_id")				// if author_id is given
	sortPara := r.URL.Query().Get("sort")					// if sort is given
	desc := false
	if sortPara == "desc" {
		desc = true
	}

	// fetch the data according to the paramater in url
	var chirpList []database.Chirp
	var err error

	if authorId == "" {
		chirpList, err = apicfg.DB.GetAllChirps(r.Context())					
		if err != nil {
			fmt.Printf("Error getting all chirps from db: %s\n", err)
			msg := "Couldn't get all chirps!"
			respondWithError(w, 500, msg)
			return 
		}
	} else {
		// convert ID string to UUID
		var userID uuid.UUID
		userID, err = uuid.Parse(authorId)
		if err != nil {
			respondWithError(w, 400, "Invalid ID")
			return
		}
		
		// get list of chirps created by this userID
		chirpList, err = apicfg.DB.GetAllChirpsFromUserID(r.Context(), userID)
		if err != nil {
			fmt.Printf("Error getting all chirps from db: %s\n", err)
			msg := "Couldn't get all chirps!"
			respondWithError(w, 500, msg)
			return
		}
	}

	// user wants desc order
	if desc {
		sort.Slice(chirpList, func(i, j int) bool {
			return chirpList[i].CreatedAt.After(chirpList[j].CreatedAt) 
		})
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