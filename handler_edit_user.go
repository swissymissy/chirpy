package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"errors"
	"database/sql"

	"github.com/swissymissy/chirpy/internal/auth"
	"github.com/swissymissy/chirpy/internal/database"
)

type newPasswordEmail struct {
	NewPassword string `json:"password"`
	NewEmail 	string `json:"email"`
}


func (apicfg *apiConfig) handlerEditUser( w http.ResponseWriter, r *http.Request) {
	// get user's access token in header
	accessToken, err := auth.GetBearerToken(r.Header) 
	if err != nil {
		fmt.Printf("Error getting token from header: %s\n", err)
		respondWithError( w, 401, "Invalid token")
		return
	}

	// check if the token is the right access token
	userID, err := auth.ValidateJWT(accessToken, apicfg.jwt_secret)
	if err != nil {
		fmt.Printf("Invalid token")
		respondWithError(w, 401, "Invalid token")
		return
	}

	//decode the request body
	decoder := json.NewDecoder(r.Body)
	var data newPasswordEmail
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Printf("Error decoding json bytes: %s\n", err)
		respondWithError(w, 500, "Something went wrong. Try again")
		return
	}

	// hash the new password
	newHashedPassword, err := auth.HashPassword(data.NewPassword)
	if err != nil {
		fmt.Printf("Error hashing new password: %s\n", err)
		respondWithError(w, 500, "Something went wrong. Try again")
		return
	}

	// update user info in db based on ID
	updatedUser, err := apicfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		ID: userID,
		Email: data.NewEmail,
		HashedPassword: newHashedPassword,
	})
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("User does not exist")
		msg := "Unauthorized"
		respondWithError(w, 401, msg)
		return
	} else if err != nil {
		fmt.Printf("Error getting row from users: %s\n", err)
		msg := "Unauthorized"
		respondWithError(w, 401, msg)
		return
	}
	respondWithJSON(w, 200 , User{
		ID: updatedUser.ID,
		Email: updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	})

}