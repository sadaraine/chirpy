package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Cleaned_body string `json:"cleaned_body"`
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":   true,
		"fornax":     true,
	}

	words := strings.Split(params.Body, " ")

	for i, word := range words {
		if badWords[strings.ToLower(word)] {
			words[i] = "****"
		}
	}

	cleaned_body := strings.Join(words, " ")

	respondWithJSON(w, http.StatusOK, returnVals{
		Cleaned_body: cleaned_body,
		Valid: true,
	})
}
