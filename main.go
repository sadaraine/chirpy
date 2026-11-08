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
}

func main() {
     const filepathRoot = "."
     const port = "8080"

     godotenv.Load()
     dbURL := os.Getenv("DB_URL")
     if dbURL == "" {
	log.Fatal("DB_URL must be set")
     }

     dbConn, err := sql.Open("postgres", dbURL)
     if err != nil {
	log.Fatal("Error opening database: %s", err)
     }
     dbQueries := database.New(dbConn)

     stripHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

     apiCfg := apiConfig{
	fileserverHits: atomic.Int32{},
        db: 		dbQueries,
     }

     mux := http.NewServeMux()
     mux.Handle("/app/", apiCfg.middlewareMetricsInc(stripHandler))
     mux.HandleFunc("GET /api/healthz", handlerReadiness)
     mux.HandleFunc("GET /admin/metrics", apiCfg.handlerCount)
     mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
     mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirps)
     
     server := &http.Server{
     	    Addr:    ":" + port,
    	    Handler: mux,
	    }

     log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
     log.Fatal(server.ListenAndServe())	
}
