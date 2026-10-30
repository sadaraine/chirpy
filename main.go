package main

import (
	"fmt"
	"net/http"
)

func readinessHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", readinessHandler)

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	server := &http.Server{
	       	Addr:    ":8080",
    		Handler: mux,
	}

	fmt.Println("Server is listening on :8080...")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}