package handler

import (
	"encoding/json"
	"net/http"
	"todo-list/src/models"
	"todo-list/src/stores"
	"todo-list/src/validations"
)

var todos []models.Todo

func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	todo := models.Todo{}

	err := json.NewDecoder(r.Body).Decode(&todo)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
	}

	errors := validations.ValidateTodo(&todo)
	if len(errors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)
		return
	}

	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

	newTodo, err := stores.GetStore().CreateTodo(&todo)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}
