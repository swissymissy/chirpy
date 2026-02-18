package main 

import (
	"fmt"
	"net/http"

	"github.com/swissymissy/chirpy/internal/auth"
)


func (apicfg *apiConfig) handlerRevokeToken( w http.ResponseWriter, r *http.Request) {
	// check for refresh token in header
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Printf("error getting refresh token from header: %s\n", err)
		respondWithError(w, 401, "Invalid token")
		return
	}

	// look for token in db
	refreshTokenInDb, err := apicfg.DB.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		fmt.Printf("Error finding token in db: %s\n", err)
		respondWithError(w, 401, "Invalid token")
		return
	}

	// reset the revoke time in db
	err = apicfg.DB.UpdateRevokedToken(r.Context(), refreshTokenInDb.Token)
	if err != nil {
		fmt.Printf("Error updating token revoked time in db: %s\n", err)
		respondWithError(w, 401 ,"Something went wrong. Try again")
		return
	}
	w.WriteHeader(204)		// response with successful request
}