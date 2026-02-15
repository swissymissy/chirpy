package main 

import (
	"fmt"
	"strings"
)


func ValidateChirp(chrpmsg *chirpMsg) error{
	
	// response struct for valid case
	type returnValid struct {
		Cleaned string `json:"cleaned_body"`
	}

	// handle error case
	if len(chrpmsg.Body) > 140 {
		return fmt.Errorf("Chirp message is too long")
	} 
	
	chrpmsg.Body = cleanerString(chrpmsg.Body)
	return nil
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
