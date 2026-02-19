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
	"github.com/swissymissy/chirpy/internal/database"
)

// response format for log in
type loginResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	RefreshToken string `json:"refresh_token"`
	IsChirpyRed	 bool	`json:"is_chirpy_red"`
}

func (apicfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	
	// email and password sent from user
	type userEP struct {
		Password string `json:"password"`
		Email string `json:"email"`
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

	// check user's password with the hased_passwd in db
	match, err := auth.CheckPasswordHash(userep.Password, userInfo.HashedPassword)
	if err != nil {
		fmt.Printf("Error checking password: %s\n", err)
		msg := "Incorrect email or password"
		respondWithError(w, 401, msg)
		return
	}
	
	if !match {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	// create a new access token for user
	token, err := auth.MakeJWT(userInfo.ID, apicfg.jwt_secret)
	if err != nil {
		fmt.Printf("Error making new access token: %s\n", err)
		respondWithError(w, 400, "Something went wrong")
		return
	}

	// create a new refresh token for user
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		fmt.Printf("Error making new refresh token: %s\n", err)
		respondWithError(w, 400, "Something went wrong")
		return
	}

	// store refresh token in db
	createRefreshTokenParams := database.CreateRefreshTokenParams {
		Token: refreshToken,
		UserID: userInfo.ID,
	}

	_, err = apicfg.DB.CreateRefreshToken(r.Context(), createRefreshTokenParams)
	if err != nil {
		fmt.Printf("Error storing refresh token to db: %s\n", err)
		respondWithError(w, 400, "Something went wrong")
		return
	}

	fmt.Println("User has logged in")
	respondWithJSON(w, 200, loginResponse{
		ID: userInfo.ID,
		CreatedAt: userInfo.CreatedAt,
		UpdatedAt: userInfo.UpdatedAt,
		Email: userInfo.Email,
		Token: token,
		RefreshToken: refreshToken,
		IsChirpyRed: userInfo.IsChirpyRed,
	})
}