package stores

import (
	"database/sql"
	"todo-list/src/models"
)

type Store interface {
	GetTodos() ([]*models.Todo, error)
	CreateTodo(todo *models.Todo) (*models.Todo, error)
	UpdateTodo(todo *models.Todo, ID int) (*models.Todo, error)
}

type DbStore struct {
	DB *sql.DB
}

var store Store

func GetStore() Store {
	return store
}

func (store *DbStore) CreateTodo(todo *models.Todo) (*models.Todo, error) {

	row := store.DB.QueryRow("INSERT INTO todos(task_name, completed, due_date) VALUES ($1, $2, $3) RETURNING task_name, completed, due_date, created_at, updated_at", todo.TaskName, todo.Completed, todo.DueDate)

	lastInsertedTodo := &models.Todo{}

	err := row.Scan(&lastInsertedTodo.TaskName, &lastInsertedTodo.Completed, &lastInsertedTodo.DueDate, &lastInsertedTodo.CreatedAt, &lastInsertedTodo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return lastInsertedTodo, nil
}

func (store *DbStore) GetTodos() ([]*models.Todo, error) {
	rows, err := store.DB.Query("Select task_name, completed, due_date FROM todos")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	todos := []*models.Todo{}
	for rows.Next() {
		todo := &models.Todo{}
		if err := rows.Scan(&todo.TaskName, &todo.Completed, &todo.DueDate); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (store *DbStore) UpdateTodo(todo *models.Todo, ID int) (*models.Todo, error) {
	row := store.DB.QueryRow("UPDATE todos SET task_name=$1, completed=$2, due_date=$3 WHERE id=$4 RETURNING task_name, completed, due_date, created_at, updated_at", todo.TaskName, todo.Completed, todo.DueDate, ID)

	updatedTodo := &models.Todo{}
	err := row.Scan(&updatedTodo.TaskName, &updatedTodo.Completed, &updatedTodo.DueDate, &updatedTodo.CreatedAt, &updatedTodo.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return updatedTodo, nil
}

func InitStore(s Store) {
	store = s
}
