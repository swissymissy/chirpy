package main 

import (
	"net/http"
)


// custom handler
// register handler func to repsonse for readiness of the server, endpoint /healthz
 
func handlehealthz(w http.ResponseWriter, req *http.Request ) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}