package main

import (
	"database/sql"
	"log"
	"net/http"
	"todo-list/src/handler"
	"todo-list/src/stores"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func routes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/todos", handler.GetTodos).Methods("GET")
	r.HandleFunc("/todos", handler.CreateTodoHandler).Methods("POST")
	r.HandleFunc("/todos/{id}", handler.UpdateTodoHandler).Methods("DELETE")
	r.HandleFunc("/users", handler.CreateUserHandler).Methods("POST")

	return r
}

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Web started!!")
// }

func main() {
	connStirng := "host=localhost port=5432 user=postgres password=secret database=todos sslmode=disable"
	db, err := sql.Open("postgres", connStirng)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	stores.InitStore(&stores.DbStore{DB: db})
	r := routes()
	log.Fatal(http.ListenAndServe(":8080", r))
}
