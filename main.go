package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

// struct to keep track of the number of requests received
type apiConfig struct {
	fileserverHits atomic.Int32
}

// middleware, increment received request
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) (http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1) 	// increment for each call
		next.ServeHTTP(w, req)
	})
}

// print amount of received hits in response
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	numHits := cfg.fileserverHits.Load()	// get number if hits from atomic.Int32
	line := fmt.Sprintf("Hits: %d", numHits)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(line))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(200)
}

func main() {
	// new apiConfig
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// create new server mux
	serverMux := http.NewServeMux()

	// create new server
	newServer := http.Server{
		Addr: ":8080",
		Handler: serverMux,
	}

	// create handler
	// func FileServer(root FileSystem) Handler
	handler := http.FileServer(http.Dir("."))

	// lil helper in case user type "/app" in the path, this will redirect them tot he right location of "/app/"
	serverMux.HandleFunc("/app",
		func(w http.ResponseWriter, req *http.Request) {
			http.Redirect(w, req, "/app/", http.StatusMovedPermanently)
		},
	)

	serverMux.HandleFunc("GET /healthz", handlehealthz)	// server's readiness check
	serverMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app",handler)))
	serverMux.HandleFunc("GET /metrics", apiCfg.handlerMetrics)
	serverMux.HandleFunc("POST /reset", apiCfg.handlerReset)

	// start the server
	err := newServer.ListenAndServe()
	if err != nil {
		fmt.Println("error listening and serve")
		return
	}
	return 
}
