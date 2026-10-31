package main

import (
	"log"
	"net/http"
)

type apiConfig struct {
     fileserverHits atomic.Int32
}

func main() {
     const filepathRoot = "."
     const port = "8080"
     stripHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

     mux := http.NewServeMux()
     mux.Handle("/app/", apiCfg.middlewareMetricsInc(stripHandler))
     mux.HandleFunc("/healthz", handlerReadiness)
     mux.HandleFunc("/metrics", apiCfg.handlerCount())
     mux.HandleFunc("/reset", apiCfg.handlerReset())
     
     server := &http.Server{
     	    Addr:    ":" + port,
    	    Handler: mux,
	    }

     log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
     log.Fatal(server.ListenAndServe())	
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          cfg.fileserverHits.Add(1)
     	  next.ServeHTTP(w,r)
     })
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
     w.Header().Add("Content-Type", "text/plain; charset=utf-8")
     w.WriteHeader(http.StatusOK)
     w.Write([]byte(fmt.Sprintf("Count Reset")
     cfg.fileserverHits.Store(0)
}

func (cfg *apiConfig) handlerCount(w http.ResponseWriter, r *http.Request) {
     w.Header().Add("Content-Type", "text/plain; charset=utf-8")
     w.WriteHeader(http.StatusOK)
     w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}