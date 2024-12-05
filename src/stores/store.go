package stores

import (
	"database/sql"
	"todo-list/src/models"
)

type Store interface {
	GetTodos(userID int) ([]*models.Todo, error)
	CreateTodo(todo *models.Todo, userID int) (*models.Todo, error)
	UpdateTodo(todo *models.Todo, todoID int) (*models.Todo, error)
	GetTodo(todoID int, userID int) error
	DeleteTodo(ID int) error
	CreateUser(user *models.User) (*models.User, error)
	GetUser(user *models.User) (*models.User, error)
}

type DbStore struct {
	DB *sql.DB
}

var store Store

func GetStore() Store {
	return store
}

func (store *DbStore) CreateTodo(todo *models.Todo, userID int) (*models.Todo, error) {
	transaction, err := store.DB.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			transaction.Rollback()
		}
	}()

	lastInsertedTodo := &models.Todo{}
	err = transaction.QueryRow("INSERT INTO todos(task_name, completed, due_date) VALUES ($1, $2, $3) RETURNING id, task_name, completed, due_date, created_at, updated_at", todo.TaskName, todo.Completed, todo.DueDate).Scan(&lastInsertedTodo.ID, &lastInsertedTodo.TaskName, &lastInsertedTodo.Completed, &lastInsertedTodo.DueDate, &lastInsertedTodo.CreatedAt, &lastInsertedTodo.UpdatedAt)

	if err != nil {
		return nil, err
	}
	_, err = transaction.Exec("INSERT INTO users_todos (user_id, todo_id) VALUES ($1, $2)", userID, lastInsertedTodo.ID)
	if err != nil {
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	return lastInsertedTodo, nil
}

func (store *DbStore) GetTodos(userID int) ([]*models.Todo, error) {
	transaction, err := store.DB.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			transaction.Rollback()
		} else {
			transaction.Commit()
		}
	}()
	rows, err := transaction.Query("SELECT t.* from todos t JOIN users_todos ut ON t.id = ut.todo_id WHERE ut.user_id = $1", userID)

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

func (store *DbStore) GetTodo(todoID int, userID int) error {
	transaction, err := store.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			transaction.Rollback()
		}
	}()

	todo := &models.Todo{}
	err = transaction.QueryRow("SELECT t.* from todos t JOIN users_todos ut ON t.id = ut.todo_id WHERE ut.user_id = $1 AND t.id = $2", userID, todoID).Scan(&todo.ID, &todo.TaskName, &todo.Completed, &todo.DueDate, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (store *DbStore) UpdateTodo(todo *models.Todo, todoID int) (*models.Todo, error) {
	updatedTodo := &models.Todo{}
	err := store.DB.QueryRow("UPDATE todos SET task_name=$1, completed=$2, due_date=$3 WHERE id=$4 RETURNING id, task_name, completed, due_date, created_at, updated_at", todo.TaskName, todo.Completed, todo.DueDate, todoID).Scan(&updatedTodo.ID, &updatedTodo.TaskName, &updatedTodo.Completed, &updatedTodo.DueDate, &updatedTodo.CreatedAt, &updatedTodo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return updatedTodo, nil
}

func (store *DbStore) DeleteTodo(ID int) error {
	_, err := store.DB.Exec("DELETE FROM todos WHERE id=$1", ID)
	if err != nil {
		return err
	}
	return nil
}

func (store *DbStore) CreateUser(user *models.User) (*models.User, error) {
	row := store.DB.QueryRow("INSERT INTO users(username, email, password) VALUES ($1, $2, $3) RETURNING username, email", user.UserName, user.Email, user.Password)
	lastInsertedUser := &models.User{}
	err := row.Scan(&lastInsertedUser.UserName, &lastInsertedUser.Email)
	if err != nil {
		return nil, err
	}
	return lastInsertedUser, nil
}

func (store *DbStore) GetUser(user *models.User) (*models.User, error) {
	row := store.DB.QueryRow("SELECT id, username, email, password FROM users WHERE email=$1 AND password=$2", user.Email, user.Password)
	userData := &models.User{}
	err := row.Scan(&userData.ID, &userData.UserName, &userData.Email, &userData.Password)
	if err != nil {
		return nil, err
	}
	return userData, nil
}

func InitStore(s Store) {
	store = s
}
