package main

import (
	"net/http"
	"time"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/swissymissy/chirpy/internal/auth"
	"github.com/swissymissy/chirpy/internal/database"
)


type User struct {
	ID 			uuid.UUID 	`json:"id"`
	CreatedAt	time.Time 	`json:"created_at"`
	UpdatedAt	time.Time 	`json:"updated_at"`
	Email		string 		`json:"email"`
}

func (apicfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	
	type reqEmailPwd struct {
		Password string `json:"password"`
		Email	string 	`json:"email"`
	}

	// decode body req into json bytes
	decoder := json.NewDecoder(r.Body)
	var rE reqEmailPwd
	err := decoder.Decode(&rE)
	if err != nil {
		fmt.Printf("Error decoding body request: %s", err)
		msg := "Something went wrong"
		respondWithError(w, 500, msg )
		return
	}

	// hash password
	pwd := rE.Password
	hashed_pwd, err := auth.HashPassword(pwd)
	if err != nil {
		fmt.Print("Error hashing password: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	// create user params
	userParams := database.CreateUserParams{
		Email: rE.Email,
		HashedPassword: hashed_pwd,
	}
	// inserting user email and hashed password to db
	newUsr, err := apicfg.DB.CreateUser(r.Context(), userParams)
	if err != nil {
		fmt.Printf("Error creating new user in db: %s", err)
		msg := "Something went wrong"
		respondWithError(w, 500, msg )
		return
	}

	// encode user struct into json
	resUsr := User{
		ID: newUsr.ID,
		CreatedAt: newUsr.CreatedAt,
		UpdatedAt: newUsr.UpdatedAt,
		Email: newUsr.Email,
	}
	respondWithJSON(w, 201, resUsr)
	return
} 