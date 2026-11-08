package main

import (
	"log"
	"net/http"
	"os"
	"database/sql"
	"sync/atomic"
	"time"
	"github.com/lpernett/godotenv"
	"github.com/sadaraine/chirpy/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)



type apiConfig struct {
     fileserverHits	atomic.Int32
     db			*database.Queries
     platform		string
}

type User struct {
	Id	uuid.UUID	`json:"id"`
	Created_At	time.Time	`json:"created_at"`
	Updated_At	time.Time	`json:"updated_at"`
	Email	string	`json:"email"`
}

func main() {
     const filepathRoot = "."
     const port = "8080"

     godotenv.Load()
     dbURL := os.Getenv("DB_URL")
     platform := os.Getenv("PLATFORM")
     if dbURL == "" {
	log.Fatal("DB_URL must be set")
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
     }

     mux := http.NewServeMux()
     mux.Handle("/app/", apiCfg.middlewareMetricsInc(stripHandler))
     mux.HandleFunc("GET /api/healthz", handlerReadiness)
     mux.HandleFunc("GET /admin/metrics", apiCfg.handlerCount)
     mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
     mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
     mux.HandleFunc("POST /api/users", apiCfg.handlerUsers)
     
     server := &http.Server{
     	    Addr:    ":" + port,
    	    Handler: mux,
	    }

     log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
     log.Fatal(server.ListenAndServe())	
}
