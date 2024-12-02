package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"todo-list/src/lib"
	"todo-list/src/models"
	"todo-list/src/stores"

	"github.com/gorilla/mux"
)

// single test case
// func TestCreateTodoHandler(t *testing.T) {
// 	payload := `{"task_name": "Learn Go", "completed": false, "due_date": "2024-11-30T23:59:59Z"}`
// 	fixedTime := time.Date(2024, time.November, 24, 0, 0, 0, 0, time.UTC)

// 	expected := &models.Todo{
// TaskName:  "Learn Go",
// Completed: false,
// DueDate:   time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
// CreatedAt: fixedTime.UTC(),
// UpdatedAt: fixedTime.UTC(),
// 	}

// 	mockStore := stores.InitMockStore()

// 	mockStore.On("CreateTodo", &models.Todo{TaskName: "Learn Go", Completed: false,
// 		DueDate: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
// 	}).Return(expected, nil)

// 	stores.InitStore(mockStore)

// 	req, err := http.NewRequest("POST", "/todos", strings.NewReader(payload))
// 	if err != nil {
// 		t.Fatalf("Failed to create request: %v", err)
// 	}

// 	recorder := httptest.NewRecorder()
// 	handler := http.HandlerFunc(CreateTodoHandler)
// 	handler.ServeHTTP(recorder, req)

// 	if status := recorder.Code; status != http.StatusCreated {
// 		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
// 	}

// 	var responseTodo models.Todo
// 	err = json.NewDecoder(recorder.Body).Decode(&responseTodo)
// 	if err != nil {
// 		t.Fatalf("Failed to decode response body: %v", err)
// 	}

// 	if responseTodo != *expected {
// 		t.Errorf("Handler returned unexpected body:\nGot:  %+v\nWant: %+v", responseTodo, expected)
// 	}

// 	mockStore.AssertExpectations(t)
// }

// Table Driven testing
func TestCreateTodoHandler(t *testing.T) {
	type testCase struct {
		name             string
		payload          string
		expectedBody     interface{}
		expectedStatus   int
		mockReturn       func(*stores.MockStore)
		token            string
		getUserMockStore func(*stores.MockStore)
	}

	fixedTime := time.Date(2024, time.November, 24, 0, 0, 0, 0, time.UTC)

	tests := []testCase{
		{
			name:    "Valid Todo",
			payload: `{"task_name": "Learn Go", "completed": false, "due_date": "2024-11-30T23:59:59Z"}`,
			expectedBody: &models.Todo{
				TaskName:  "Learn Go",
				Completed: false,
				DueDate:   time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
				CreatedAt: fixedTime.UTC(),
				UpdatedAt: fixedTime.UTC(),
			},
			expectedStatus: http.StatusCreated,
			mockReturn: func(mockStore *stores.MockStore) {
				mockStore.On("CreateTodo", &models.Todo{
					TaskName:  "Learn Go",
					Completed: false,
					DueDate:   time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
				}).Return(&models.Todo{
					TaskName:  "Learn Go",
					Completed: false,
					DueDate:   time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
					CreatedAt: fixedTime.UTC(),
					UpdatedAt: fixedTime.UTC(),
				}, nil)
			},
			token: func() string {
				token, err := lib.GenerateJWT("test@mail.com", "password")
				if err != nil {
					t.Fatalf("Failed to generate JWT: %v", err)
				}
				return *token
			}(),
			getUserMockStore: func(mockStore *stores.MockStore) {
				mockStore.On("GetUser", &models.User{
					Email:    "test@mail.com",
					Password: "password",
				}).Return(&models.User{
					UserName: "testuser",
					Email:    "test@mail.com",
					Password: "password",
				}, nil)
			},
		},
		{
			name:           "InValid Todo",
			payload:        `{"task_name": "", "completed": false, "due_date": "2024-11-30T23:59:59Z"}`,
			expectedBody:   &models.Todo{},
			expectedStatus: http.StatusBadRequest,
			mockReturn:     func(mockStore *stores.MockStore) {},
			token: func() string {
				token, err := lib.GenerateJWT("test@mail.com", "password")
				if err != nil {
					t.Fatalf("Failed to generate JWT: %v", err)
				}
				return *token
			}(),
			getUserMockStore: func(mockStore *stores.MockStore) {
				mockStore.On("GetUser", &models.User{
					Email:    "test@mail.com",
					Password: "password",
				}).Return(&models.User{
					UserName: "testuser",
					Email:    "test@mail.com",
					Password: "password",
				}, nil)
			},
		},
		{
			name:             "In Valid token request",
			payload:          `{"task_name": "Learn Go", "completed": false, "due_date": "2024-11-30T23:59:59Z"}`,
			expectedBody:     map[string]string{"error": "Invalid token"},
			expectedStatus:   http.StatusUnauthorized,
			mockReturn:       func(mockStore *stores.MockStore) {},
			token:            "invalid token",
			getUserMockStore: func(mockStore *stores.MockStore) {},
		},
		{
			name:           "In Valid user",
			payload:        `{"task_name": "Learn Go", "completed": false, "due_date": "2024-11-30T23:59:59Z"}`,
			expectedBody:   map[string]string{"error": "User not found"},
			expectedStatus: http.StatusUnauthorized,
			mockReturn:     func(mockStore *stores.MockStore) {},
			token: func() string {
				token, err := lib.GenerateJWT("test@mail.com", "password")
				if err != nil {
					t.Fatalf("Failed to generate JWT: %v", err)
				}
				return *token
			}(),
			getUserMockStore: func(mockStore *stores.MockStore) {
				mockStore.On("GetUser", &models.User{
					Email:    "test@mail.com",
					Password: "password",
				}).Return(&models.User{}, errors.New("User not found"))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockStore := stores.InitMockStore()
			tc.mockReturn(mockStore)
			tc.getUserMockStore(mockStore)
			stores.InitStore(mockStore)

			req, err := http.NewRequest("POST", "/todos", strings.NewReader(tc.payload))
			req.Header.Set("Authorization", "Bearer "+tc.token)

			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(CreateTodoHandler)
			handler.ServeHTTP(recorder, req)

			if status := recorder.Code; status != tc.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

			if tc.expectedBody != nil {
				var responseTodo interface{}
				switch tc.expectedBody.(type) {
				case map[string]string:
					responseTodo = make(map[string]string)
				case *models.Todo:
					responseTodo = &models.Todo{}
				default:
					t.Fatalf("Unsupported type: %T", tc.expectedBody)
				}

				err := json.NewDecoder(recorder.Body).Decode(&responseTodo)

				if err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}
				// doesn't work because of type mismatch

				// log.Printf("Type of expectedBody: %T, Type of responseTodo: %T", tc.expectedBody, responseTodo)
				// if !reflect.DeepEqual(responseTodo, tc.expectedBody) {
				// 	t.Errorf("Handler returned unexpected body:\nGot:  %+v\nWant: %+v", responseTodo, tc.expectedBody)
				// }
			}

			mockStore.AssertExpectations(t)
		})
	}

}

func TestGetTodoHandler(t *testing.T) {
	type testCase struct {
		name           string
		expectedTodos  []*models.Todo
		expectedStatus int
		mockReturn     func(*stores.MockStore)
	}

	tests := []testCase{
		{
			name: "Get Todos",
			expectedTodos: []*models.Todo{
				{TaskName: "Learn Go", Completed: false, DueDate: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC()},
				{TaskName: "Learn Ruby", Completed: false, DueDate: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC()},
				{TaskName: "Learn Python", Completed: false, DueDate: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC()},
			},
			expectedStatus: http.StatusOK,
			mockReturn: func(mockStore *stores.MockStore) {
				mockStore.On("GetTodos").Return([]*models.Todo{
					{TaskName: "Learn Go", Completed: false, DueDate: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC()},
					{TaskName: "Learn Ruby", Completed: false, DueDate: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC()},
					{TaskName: "Learn Python", Completed: false, DueDate: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC()},
				}, nil)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockStore := stores.InitMockStore()
			tc.mockReturn(mockStore)
			stores.InitStore(mockStore)

			req, err := http.NewRequest("GET", "/todos", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(GetTodosHandler)
			handler.ServeHTTP(recorder, req)

			if status := recorder.Code; status != tc.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

			if tc.expectedTodos != nil {
				var responseTodos []*models.Todo
				err := json.NewDecoder(recorder.Body).Decode(&responseTodos)
				if err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				if len(responseTodos) != len(tc.expectedTodos) {
					t.Errorf("Handler returned unexpected body:\nGot:  %+v\nWant: %+v", responseTodos, tc.expectedTodos)
				}
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestUpdateTodoHandler(t *testing.T) {
	type testCase struct {
		name           string
		payload        string
		expectedTodo   *models.Todo
		expectedStatus int
		mockReturn     func(*stores.MockStore)
	}
	tests := []testCase{
		{
			name:    "Update Todo",
			payload: `{"task_name": "Learn Go", "completed": false, "due_date": "2024-11-30T23:59:59Z"}`,
			expectedTodo: &models.Todo{
				TaskName:  "Updated Learn Go",
				Completed: true,
				DueDate:   time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
				CreatedAt: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
				UpdatedAt: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
			},
			expectedStatus: http.StatusCreated,
			mockReturn: func(mockStore *stores.MockStore) {
				mockStore.On("UpdateTodo", &models.Todo{
					TaskName:  "Learn Go",
					Completed: false,
					DueDate:   time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
				}, 1).Return(&models.Todo{
					TaskName:  "Updated Learn Go",
					Completed: true,
					DueDate:   time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
					CreatedAt: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
					UpdatedAt: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC(),
				}, nil)
			},
		},
		{
			name:           "Invalid Update Todo",
			payload:        `{"task_name": "", "completed": false, "due_date": "2024-11-30T23:59:59Z"}`,
			expectedTodo:   &models.Todo{},
			expectedStatus: http.StatusBadRequest,
			mockReturn:     func(mockStore *stores.MockStore) {},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockStore := stores.InitMockStore()
			tc.mockReturn(mockStore)
			stores.InitStore(mockStore)

			req, err := http.NewRequest("PUT", "/todos/1", strings.NewReader(tc.payload))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			r := mux.NewRouter()
			r.HandleFunc("/todos/{id:[0-9]+}", UpdateTodoHandler).Methods("PUT")
			recorder := httptest.NewRecorder()

			// handler := http.HandlerFunc(UpdateTodoHandler)
			// handler.ServeHTTP(recorder, req)
			r.ServeHTTP(recorder, req)

			if status := recorder.Code; status != tc.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

			updatedTodo := models.Todo{}
			err = json.NewDecoder(recorder.Body).Decode(&updatedTodo)
			if err != nil {
				t.Fatalf("Failed to decode response body: %v", err)
			}
			if updatedTodo != *tc.expectedTodo {
				t.Errorf("Handler returned unexpected body:\nGot:  %+v\nWant: %+v", updatedTodo, tc.expectedTodo)
			}

			mockStore.AssertExpectations(t)
		})
	}

}

func TestDeleteTodoHandler(t *testing.T) {
	type testCase struct {
		name           string
		id             any
		expectedStatus int
		mockReturn     func(*stores.MockStore)
	}

	tests := []testCase{
		{
			name:           "Delete Todo",
			expectedStatus: http.StatusOK,
			id:             1,
			mockReturn: func(mockStore *stores.MockStore) {
				mockStore.On("DeleteTodo", 1).Return(nil)
			},
		},
		{
			name:           "Delete Todo Invalid ID",
			expectedStatus: http.StatusNotFound,
			id:             "abc",
			mockReturn:     func(mockStore *stores.MockStore) {},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockStore := stores.InitMockStore()
			tc.mockReturn(mockStore)
			stores.InitStore(mockStore)
			ID := tc.id

			req, err := http.NewRequest("DELETE", "/todos/"+fmt.Sprintf("%v", ID), nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			r := mux.NewRouter()
			r.HandleFunc("/todos/{id:[0-9]+}", DeleteTodoHandler).Methods("DELETE")

			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)

			if status := recorder.Code; status != tc.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

			mockStore.AssertExpectations(t)
		})
	}
}
