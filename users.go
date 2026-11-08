package main

import (
	"log"
	"encoding/json"
	"net/http"
	)

type CreateUserRequest struct {
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var req CreateUserRequest
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), req.Email)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	resp := User{
		Id:	user.ID,
		Created_At:	user.CreatedAt,
		Updated_At:	user.UpdatedAt,
		Email:	user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
