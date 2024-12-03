package handler

import (
	"encoding/json"
	"fmt"
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
		utility.WriteJsonData(w, map[string]string{"error": "Invalid request payload"}, http.StatusBadRequest)
		return
	}

	jwtToken, err := utility.ExtractTokenFromHeader(r)
	if err != nil {
		utility.WriteJsonData(w, map[string]string{"error": "Invalid token"}, http.StatusUnauthorized)
		return
	}
	claims, err := lib.ValidateJWT(jwtToken)
	if err != nil {
		utility.WriteJsonData(w, map[string]string{"error": "Invalid token"}, http.StatusUnauthorized)
		return
	}
	user := &models.User{}
	user.Email = claims["email"].(string)
	user.Password = claims["password"].(string)
	user, err = stores.GetStore().GetUser(user)

	if err != nil {
		utility.WriteJsonData(w, map[string]string{"error": "User not found"}, http.StatusUnauthorized)
		return
	}

	errors := validations.ValidateTodo(&todo)
	if len(errors) > 0 {
		utility.WriteJsonData(w, errors, http.StatusBadRequest)
		return
	}

	newTodo, err := stores.GetStore().CreateTodo(&todo, user.ID)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	utility.WriteJsonData(w, newTodo, http.StatusCreated)
}

func GetTodosHandler(w http.ResponseWriter, r *http.Request) {
	jwtToken, err := utility.ExtractTokenFromHeader(r)
	if err != nil {
		utility.WriteJsonData(w, map[string]string{"error": "Invalid token"}, http.StatusUnauthorized)
		return
	}
	claims, err := lib.ValidateJWT(jwtToken)
	if err != nil {
		utility.WriteJsonData(w, map[string]string{"error": "Invalid token"}, http.StatusUnauthorized)
		return
	}
	user := &models.User{}
	user.Email = claims["email"].(string)
	user.Password = claims["password"].(string)
	user, err = stores.GetStore().GetUser(user)

	if err != nil {
		utility.WriteJsonData(w, map[string]string{"error": "User not found"}, http.StatusUnauthorized)
		return
	}

	todos, err := stores.GetStore().GetTodos(user.ID)
	if err != nil {
		utility.WriteJsonData(w, map[string]string{"error": fmt.Sprintf("Can not get todos\n%v", err)}, http.StatusForbidden)
		return
	}

	utility.WriteJsonData(w, todos, http.StatusOK)
}

func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	todo := models.Todo{}
	vars := mux.Vars(r)
	todoID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Can not convert id to int", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		utility.WriteJsonData(w, map[string]string{"error": "Invalid request payload"}, http.StatusBadRequest)
		return
	}

	jwtToken, err := utility.ExtractTokenFromHeader(r)
	if err != nil {
		utility.WriteJsonData(w, map[string]string{"error": "Invalid token"}, http.StatusUnauthorized)
		return
	}
	claims, err := lib.ValidateJWT(jwtToken)
	if err != nil {
		utility.WriteJsonData(w, map[string]string{"error": "Invalid token"}, http.StatusUnauthorized)
		return
	}
	user := &models.User{}
	user.Email = claims["email"].(string)
	user.Password = claims["password"].(string)
	user, err = stores.GetStore().GetUser(user)

	if err != nil {
		utility.WriteJsonData(w, map[string]string{"error": "User not found"}, http.StatusUnauthorized)
		return
	}

	err = stores.GetStore().GetTodo(todoID, user.ID)
	if err != nil {
		utility.WriteJsonData(w, map[string]string{"error": "You are not authorized to update this todo"}, http.StatusNotFound)
		return
	}

	errors := validations.ValidateTodo(&todo)
	if len(errors) > 0 {
		utility.WriteJsonData(w, errors, http.StatusBadRequest)
		return
	}

	updatedTodo, err := stores.GetStore().UpdateTodo(&todo, todoID)
	if err != nil {
		utility.WriteJsonData(w, map[string]string{"error": fmt.Sprintf("Can not update todo\n%v", err)}, http.StatusForbidden)
		return
	}

	utility.WriteJsonData(w, updatedTodo, http.StatusCreated)
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
