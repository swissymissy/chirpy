package main

import(
	"net/http"
	"log"
	"encoding/json"
)

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