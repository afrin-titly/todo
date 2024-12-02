package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"todo-list/src/lib"
	"todo-list/src/models"
	"todo-list/src/stores"
	"todo-list/src/validations"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	errors := validations.ValidateUser(&user)

	if len(errors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)
		return
	}
	newUser, err := stores.GetStore().CreateUser(&user)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	req := models.User{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := stores.GetStore().GetUser(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Wrong email or password"})
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	token, err := lib.GenerateJWT(user.Email, user.Password)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": *token})
}
