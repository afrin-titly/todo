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
		mockSetup    func(todoInput *models.Todo, userID int)
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
			mockSetup: func(todoInput *models.Todo, userID int) {
				mock.ExpectQuery("INSERT INTO todos").WithArgs(todoInput.TaskName, todoInput.Completed, todoInput.DueDate).WillReturnRows(sqlmock.NewRows([]string{"id", "task_name", "completed", "due_date", "created_at", "updated_at"}).AddRow(1, todoInput.TaskName, todoInput.Completed, todoInput.DueDate, todoInput.DueDate, todoInput.DueDate))

				mock.ExpectExec("INSERT INTO users_todos").WithArgs(userID, 1).WillReturnResult(sqlmock.NewResult(1, 1))
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
			expectedTodo: &models.Todo{
				ID:        1,
				TaskName:  "test task",
				Completed: false,
				DueDate:   time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
				CreatedAt: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
				UpdatedAt: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
			},
			userID: 1,
			mockSetup: func(todoInput *models.Todo, userID int) {
				mock.ExpectQuery("INSERT INTO todos").WithArgs(todoInput.TaskName, todoInput.Completed, todoInput.DueDate).WillReturnError(fmt.Errorf("error inserting into todos"))
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
			mockSetup: func(todoInput *models.Todo, userID int) {
				mock.ExpectQuery("INSERT INTO todos").WithArgs(todoInput.TaskName, todoInput.Completed, todoInput.DueDate).WillReturnRows(sqlmock.NewRows([]string{"id", "task_name", "completed", "due_date", "created_at", "updated_at"}).AddRow(1, todoInput.TaskName, todoInput.Completed, todoInput.DueDate, todoInput.DueDate, todoInput.DueDate))

				mock.ExpectExec("INSERT INTO users_todos").WithArgs(userID, 1).WillReturnError(fmt.Errorf("some db error"))
			},
			shouldError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock.ExpectBegin()
			tc.mockSetup(tc.todoInput, tc.userID)
			mock.ExpectCommit()
			newTodo, err := store.CreateTodo(tc.todoInput, tc.userID)
			if tc.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedTodo, newTodo)
			}

		})
	}
}
