package main 

import (
	"encoding/json"
	"net/http"
	"log"
	"strings"
)


func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	// response struct for valid case
	type returnValid struct {
		Cleaned string `json:"cleaned_body"`
	}

	// decode body req into json bytes
	decoder := json.NewDecoder(req.Body)
	var params parameters
	err := decoder.Decode(&params)	// write to params after decoding
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		msg := "Something went wrong"
		respondWithError(w, 500, msg )
		return
	}

	// handle error case
	if len(params.Body) > 140 {
		msg := "Chirp is too long"
		respondWithError(w, 400, msg)
		return
	} 
	
	cleaned_res := cleanerString(params.Body)
	res := returnValid{
		Cleaned: cleaned_res,
	}
	respondWithJSON(w, 200, res)
}

// helper funcs
func respondWithError( w http.ResponseWriter, code int, msg string) {
	type returnWithErr struct {
		Error string `json:"error"`
	}

	res := returnWithErr{
		Error: msg,
	}
	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error encoding msg to json: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithJSON( w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	
	//encode to json bytes
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error encoding payload to json: %s", err)
		return
	}
	w.Write(data)
}

// cleaner functions
func cleanerString(msg string) string {
	split_msg := strings.Fields(msg) 
	// set of bad words
	bad := map[string]struct{}{
		"kerfuffle":{},
		"sharbert":{},
		"fornax":{},
	}

	for i := range split_msg{
		word := strings.ToLower(split_msg[i])
		if _, ok := bad[word]; ok {
			split_msg[i] = "****"
		}
	}
	msg = strings.Join(split_msg, " ")
	return msg
}
