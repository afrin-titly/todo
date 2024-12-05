package stores

import (
	"fmt"
	"testing"
	"time"
	"todo-list/src/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTodo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	store := &DbStore{DB: db}

	type testCase struct {
		name         string
		todoInput    *models.Todo
		expectedTodo *models.Todo
		userID       int
		mockSetup    func(todoInput *models.Todo, userID int, expectedTodo *models.Todo)
		shouldError  bool
	}

	tests := []testCase{
		{
			name: "Successful todo creation",
			todoInput: &models.Todo{
				TaskName:  "test task",
				Completed: false,
				DueDate:   time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
			},
			expectedTodo: &models.Todo{
				ID:        1,
				TaskName:  "test task",
				Completed: false,
				DueDate:   time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
				CreatedAt: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
				UpdatedAt: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
			},
			userID: 1,
			mockSetup: func(todoInput *models.Todo, userID int, expectedTodo *models.Todo) {
				mock.ExpectQuery("INSERT INTO todos").WithArgs(todoInput.TaskName, todoInput.Completed, todoInput.DueDate).WillReturnRows(sqlmock.NewRows([]string{"id", "task_name", "completed", "due_date", "created_at", "updated_at"}).AddRow(1, expectedTodo.TaskName, expectedTodo.Completed, expectedTodo.DueDate, expectedTodo.DueDate, expectedTodo.DueDate))

				mock.ExpectExec("INSERT INTO users_todos").WithArgs(userID, 1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			shouldError: false,
		},
		{
			name: "Error in INSERT INTO todos",
			todoInput: &models.Todo{
				TaskName:  "test task",
				Completed: false,
				DueDate:   time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
			},
			expectedTodo: nil,
			userID:       1,
			mockSetup: func(todoInput *models.Todo, userID int, expectedTodo *models.Todo) {
				mock.ExpectQuery("INSERT INTO todos").WithArgs(todoInput.TaskName, todoInput.Completed, todoInput.DueDate).WillReturnError(fmt.Errorf("error inserting into todos"))
				mock.ExpectRollback()
			},
			shouldError: true,
		},
		{
			name: "Error in INSERT INTO users_todos",
			todoInput: &models.Todo{
				TaskName:  "test task",
				Completed: false,
				DueDate:   time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
			},
			expectedTodo: &models.Todo{
				ID:        1,
				TaskName:  "test task",
				Completed: false,
				DueDate:   time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
				CreatedAt: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
				UpdatedAt: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
			},
			userID: 1,
			mockSetup: func(todoInput *models.Todo, userID int, expectedTodo *models.Todo) {
				mock.ExpectQuery("INSERT INTO todos").WithArgs(todoInput.TaskName, todoInput.Completed, todoInput.DueDate).WillReturnRows(sqlmock.NewRows([]string{"id", "task_name", "completed", "due_date", "created_at", "updated_at"}).AddRow(1, expectedTodo.TaskName, expectedTodo.Completed, expectedTodo.DueDate, expectedTodo.DueDate, expectedTodo.DueDate))

				mock.ExpectExec("INSERT INTO users_todos").WithArgs(userID, 1).WillReturnError(fmt.Errorf("some db error"))
				mock.ExpectRollback()

			},
			shouldError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock.ExpectBegin()
			tc.mockSetup(tc.todoInput, tc.userID, tc.expectedTodo)
			newTodo, err := store.CreateTodo(tc.todoInput, tc.userID)
			if tc.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedTodo, newTodo)
			}
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestUpdateTodo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	store := &DbStore{DB: db}

	type testCase struct {
		name         string
		todoInput    *models.Todo
		expectedTodo *models.Todo
		todoID       int
		mockSetup    func(todoInput *models.Todo, userID int, expectedTodo *models.Todo)
		shouldError  bool
	}

	tests := []testCase{
		{
			name: "Successful update",
			todoInput: &models.Todo{
				TaskName:  "test task",
				Completed: false,
				DueDate:   time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
			},
			expectedTodo: &models.Todo{
				ID:        1,
				TaskName:  "updated test task",
				Completed: true,
				DueDate:   time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
				CreatedAt: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
				UpdatedAt: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
			},
			todoID: 1,
			mockSetup: func(todoInput *models.Todo, todoID int, expectedTodo *models.Todo) {
				mock.ExpectQuery("UPDATE todos SET task_name=\\$1, completed=\\$2, due_date=\\$3 WHERE id=\\$4 RETURNING id, task_name, completed, due_date, created_at, updated_at").WithArgs(todoInput.TaskName, todoInput.Completed, todoInput.DueDate, todoID).WillReturnRows(sqlmock.NewRows([]string{"id", "task_name", "due_date", "completed", "created_at", "updated_at"}).AddRow(expectedTodo.ID, expectedTodo.TaskName, expectedTodo.Completed, expectedTodo.DueDate, expectedTodo.CreatedAt, expectedTodo.UpdatedAt))
			},
			shouldError: false,
		},
		{
			name: "Unsuccessful update",
			todoInput: &models.Todo{
				TaskName:  "test task",
				Completed: false,
				DueDate:   time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
			},
			expectedTodo: nil,
			todoID:       1,
			mockSetup: func(todoInput *models.Todo, todoID int, expectedTodo *models.Todo) {
				mock.ExpectQuery("UPDATE todos SET task_name=\\$1, completed=\\$2, due_date=\\$3 WHERE id=\\$4 RETURNING id, task_name, completed, due_date, created_at, updated_at").WithArgs(todoInput.TaskName, todoInput.Completed, todoInput.DueDate, todoID).WillReturnError(fmt.Errorf("some db error"))
			},
			shouldError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup(tc.todoInput, tc.todoID, tc.expectedTodo)
			updatedTodo, err := store.UpdateTodo(tc.todoInput, tc.todoID)
			if tc.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedTodo, updatedTodo)
			}
		})
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	}
}

func TestGetTodos(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	store := &DbStore{DB: db}

	type testCase struct {
		name          string
		userID        int
		expectedTodos []*models.Todo
		mockSetup     func(userID int, expectedTodos []*models.Todo)
		shouldError   bool
	}

	tests := []testCase{
		{
			name:   "Successful Get todos",
			userID: 1,
			expectedTodos: []*models.Todo{
				{TaskName: "test task 1", Completed: false, DueDate: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC()},
				{TaskName: "test task 2", Completed: true, DueDate: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC()},
				{TaskName: "test task 3", Completed: false, DueDate: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC()},
			},
			mockSetup: func(userID int, expectedTodos []*models.Todo) {
				rows := sqlmock.NewRows([]string{"task_name", "completed", "due_date"})
				for _, todo := range expectedTodos {
					rows.AddRow(todo.TaskName, todo.Completed, todo.DueDate)
				}
				mock.ExpectQuery("SELECT t.* from todos t JOIN users_todos ut ON t.id = ut.todo_id WHERE ut.user_id = \\$1").WithArgs(userID).WillReturnRows(rows)
				mock.ExpectCommit()
			},
			shouldError: false,
		},
		{
			name:          "Unsuccessful Get todos",
			userID:        1,
			expectedTodos: nil,
			mockSetup: func(userID int, expectedTodos []*models.Todo) {
				mock.ExpectQuery("SELECT t.* from todos t JOIN users_todos ut ON t.id = ut.todo_id WHERE ut.user_id = \\$1").WithArgs(userID).WillReturnError(fmt.Errorf("some db error"))
				mock.ExpectRollback()
			},
			shouldError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock.ExpectBegin()
			tc.mockSetup(tc.userID, tc.expectedTodos)
			todos, err := store.GetTodos(tc.userID)
			if tc.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tc.expectedTodos), len(todos))
			}
		})
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	}
}
