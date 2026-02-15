package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"github.com/swissymissy/chirpy/internal/database"
)

// struct to keep track of the number of requests received
type apiConfig struct {
	fileserverHits atomic.Int32
	DB 	*database.Queries
	platform string 		// check if is dev
}


func main() {
	
	godotenv.Load()							// Load the .env file into enviroment variables
	dbURL := os.Getenv("DB_URL") 			// get the db_url from the environment
	platform := os.Getenv("PLATFORM")		// get PLATFORM value
	db, err := sql.Open("postgres", dbURL)	// open a connection to the database
	if err != nil {
		fmt.Println("Error connecting with database")
		return
	}
	dbQueries := database.New(db)			

	// new apiConfig
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		DB: dbQueries,
		platform: platform,
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

	serverMux.HandleFunc("GET /api/healthz", handlehealthz)	// server's readiness check
	serverMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app",handler)))
	serverMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serverMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	serverMux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	serverMux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	serverMux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	// start the server
	err = newServer.ListenAndServe()
	if err != nil {
		fmt.Println("error listening and serve")
		return
	}
	return 
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
	line := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, numHits)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(line))
}

// reset user table, only "dev" has permission
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	// check if it is a dev
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	
	err := cfg.DB.ResetUsers(req.Context())
	if err != nil {
		fmt.Printf("Error reseting users table in db: %v\n", err)
		msg := "Something went wrong"
		respondWithError(w, 500, msg)
		return
	}
	w.WriteHeader(200)
}