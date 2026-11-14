package main

import (
	"log"
	"net/http"
	"os"
	"database/sql"
	"sync/atomic"
	"github.com/lpernett/godotenv"
	"github.com/sadaraine/chirpy/internal/database"
	_ "github.com/lib/pq"
)



type apiConfig struct {
     fileserverHits	atomic.Int32
     db			*database.Queries
     platform		string
     jwtSecret		string
}

func main() {
     const filepathRoot = "."
     const port = "8080"

     godotenv.Load()
     dbURL := os.Getenv("DB_URL")
     if dbURL == "" {
	log.Fatal("DB_URL must be set")
     }

     platform := os.Getenv("PLATFORM")
     if platform == "" {
	log.Fatal("PLATFORM must be set")
     }

     jwtSecret := os.Getenv("JWT_SECRET")
     if jwtSecret == "" {
	log.Fatal("JWT_SECRET environment variable is not set")
     }

     dbConn, err := sql.Open("postgres", dbURL)
     if err != nil {
	log.Fatalf("Error opening database: %s", err)
     }
     dbQueries := database.New(dbConn)

     stripHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

     apiCfg := apiConfig{
	fileserverHits: atomic.Int32{},
        db: 		dbQueries,
	platform:	platform,
	jwtSecret:	jwtSecret,
     }

     mux := http.NewServeMux()
     mux.Handle("/app/", apiCfg.middlewareMetricsInc(stripHandler))
     mux.HandleFunc("GET /api/healthz", handlerReadiness)
     mux.HandleFunc("GET /admin/metrics", apiCfg.handlerCount)
     mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
     mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
     mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
     mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
     mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpsGet)
     mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
     
     server := &http.Server{
     	    Addr:    ":" + port,
    	    Handler: mux,
	    }

     log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
     log.Fatal(server.ListenAndServe())	
}
