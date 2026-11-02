package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
     fileserverHits atomic.Int32
}

func main() {
     const filepathRoot = "."
     const port = "8080"
     stripHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

     apiCfg := apiConfig{
	fileserverHits: atomic.Int32{},
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
