package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"todo-list/src/lib"
	"todo-list/src/models"
	"todo-list/src/stores"
	"todo-list/src/utility"
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

	jwtToken, err := utility.ExtractTokenFromHeader(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token"})
		return
	}
	claims, err := lib.ValidateJWT(jwtToken)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token"})
		return
	}

	user := &models.User{}
	user.Email = claims["email"].(string)
	user.Password = claims["password"].(string)

	user, err = stores.GetStore().GetUser(user)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		return
	}

	errors := validations.ValidateTodo(&todo)
	if len(errors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)
		return
	}

	newTodo, err := stores.GetStore().CreateTodo(&todo, user.ID)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

func GetTodosHandler(w http.ResponseWriter, r *http.Request) {
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

	newTodo, err := stores.GetStore().UpdateTodo(&todo, id)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, "Can not convert id to int", http.StatusBadRequest)
		return
	}
	err = stores.GetStore().DeleteTodo(ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"message": "Can not delete todo"})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo deleted successfully. ID: " + vars["id"]})
}
