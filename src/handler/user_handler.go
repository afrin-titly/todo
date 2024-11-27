package handler

import (
	"encoding/json"
	"net/http"
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
