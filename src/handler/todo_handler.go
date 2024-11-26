package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"todo-list/src/models"
	"todo-list/src/stores"
	"todo-list/src/validations"

	"github.com/gorilla/mux"
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

func GetTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := stores.GetStore().GetTodos()
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todos)
}

func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	todo := models.Todo{}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Can not convert id to int", http.StatusBadRequest)
		return

	}

	err = json.NewDecoder(r.Body).Decode(&todo)

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

	newTodo, err := stores.GetStore().UpdateTodo(&todo, id)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}
