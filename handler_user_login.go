package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"errors"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/swissymissy/chirpy/internal/auth"
)

type loginResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (apicfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	
	// email and password sent from user
	type userEP struct {
		Password string `json:"password"`
		Email string `json:"email"`
		ExpiresInSecond int `json:"expires_in_seconds"`
	}

	// decode response body
	decoder := json.NewDecoder(r.Body)
	var userep userEP
	err := decoder.Decode(&userep)
	if err != nil {
		fmt.Printf("Error decoding request body: %s\n", err)
		respondWithError(w, 400, "Something went wrong. Please try again")
		return
	}

	// set expire duration
	expires := time.Hour
	if userep.ExpiresInSecond > 0 {
		expires = time.Duration(userep.ExpiresInSecond) * time.Second
		if time.Duration(userep.ExpiresInSecond) > time.Hour { 
			expires = time.Hour
		} 
	}

	// get user info from db
	userInfo, err := apicfg.DB.GetUserByEmail(r.Context(), userep.Email )
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Printf("Error getting a row from users: %s\n", err)
		msg := "Incorrect email or password"
		respondWithError(w, 401, msg)
		return
	} else if err != nil {
		fmt.Printf("Error getting row from users: %s\n", err)
		msg := "Incorrect email or password"
		respondWithError(w, 401, msg)
		return
	}

	// create a new token for user
	token, err := auth.MakeJWT(userInfo.ID, apicfg.jwt_secret, expires)
	if err != nil {
		fmt.Printf("Error making new token: %s\n", err)
		respondWithError(w, 400, "Something went wrong")
		return
	}

	// compare user's input with the hased_passwd in db
	match, err := auth.CheckPasswordHash(userep.Password, userInfo.HashedPassword)
	if err != nil {
		fmt.Printf("Error checking password: %s\n", err)
		msg := "Incorrect email or password"
		respondWithError(w, 401, msg)
		return
	}

	if match {
		fmt.Println("User has logged in")
		respondWithJSON(w, 200, loginResponse{
			ID: userInfo.ID,
			CreatedAt: userInfo.CreatedAt,
			UpdatedAt: userInfo.UpdatedAt,
			Email: userInfo.Email,
			Token: token,
		})
		return
	} else {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}
}