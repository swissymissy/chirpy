package main
import (
	"log"
	"encoding/json"
	"net/http"
)

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