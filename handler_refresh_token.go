package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/swissymissy/chirpy/internal/auth"
)

type responseToken struct {
	Token string `json:"token"`
}

// function used for checking user's refresh token expired/ revoked yet. create new access token for user
func (apicfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	// get refresh token 
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Printf("error getting token from header: %s\n", err)
		respondWithError(w , 401, "Invalid token")
		return
	}

	// get refresh token from db
	refreshTokenDb, err := apicfg.DB.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		fmt.Printf("error getting user from db: %s\n", err)
		respondWithError(w , 401, "Invalid token")
		return
	}

	// check if token expires yet
	if refreshTokenDb.ExpiresAt.Before(time.Now()) {
		respondWithError(w, 401 , "Token has expired")
		return
	} 
	// check if token is revoked yet
	if refreshTokenDb.RevokedAt.Valid {
		respondWithError(w , 401, "Token has been revoked")
		return
	}

	// create new access token for user
	newAccessToken, err := auth.MakeJWT(refreshTokenDb.UserID, apicfg.jwt_secret)
	if err != nil {
		fmt.Printf("Error making new access token: %s\n", err)
		respondWithError(w , 401, "Something went wrong. Try again")
		return
	}
	respondWithJSON(w, 200 , responseToken{
		Token: newAccessToken,
	})
}
