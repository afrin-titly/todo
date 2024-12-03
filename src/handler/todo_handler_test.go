package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
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
				}, 0).Return(&models.Todo{
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
				if todo, ok := tc.expectedBody.(*models.Todo); ok {
					var decodedTodo models.Todo
					if err := json.NewDecoder(recorder.Body).Decode(&decodedTodo); err != nil {
						t.Fatalf("Failed to decode response body: %v", err)
					}

					if decodedTodo != *todo {
						t.Errorf("Handler returned unexpected body:\nGot:  %+v\nWant: %+v", decodedTodo, todos)
					}
				} else if errorBody, ok := tc.expectedBody.(map[string]string); ok {
					var decodedErrorBody map[string]string
					if err := json.NewDecoder(recorder.Body).Decode(&decodedErrorBody); err != nil {
						t.Fatalf("Failed to decode response body: %v", err)
					}

					if !reflect.DeepEqual(decodedErrorBody, errorBody) {
						t.Errorf("Handler returned unexpected body:\nGot:  %+v\nWant: %+v", decodedErrorBody, errorBody)
					}
				} else {
					t.Fatalf("Unsupported type for expectedBody: %T", tc.expectedBody)
				}
			}

			mockStore.AssertExpectations(t)
		})
	}

}

func TestGetTodosHandler(t *testing.T) {
	type testCase struct {
		name             string
		expectedBody     interface{}
		expectedStatus   int
		mockReturn       func(*stores.MockStore)
		token            string
		getUserMockStore func(*stores.MockStore)
	}

	tests := []testCase{
		{
			name: "Get Todos Success",
			expectedBody: []*models.Todo{
				{TaskName: "Learn Go", Completed: false, DueDate: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC()},
				{TaskName: "Learn Ruby", Completed: false, DueDate: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC()},
				{TaskName: "Learn Python", Completed: false, DueDate: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC()},
			},
			expectedStatus: http.StatusOK,
			mockReturn: func(mockStore *stores.MockStore) {
				mockStore.On("GetTodos", 1).Return([]*models.Todo{
					{TaskName: "Learn Go", Completed: false, DueDate: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC()},
					{TaskName: "Learn Ruby", Completed: false, DueDate: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC()},
					{TaskName: "Learn Python", Completed: false, DueDate: time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC).UTC()},
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
					ID:       1,
					UserName: "testuser",
					Email:    "test@mail.com",
					Password: "password",
				}, nil)
			},
		},
		{
			name:           "In Valid user",
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

			req, err := http.NewRequest("GET", "/todos", nil)
			req.Header.Set("Authorization", "Bearer "+tc.token)

			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(GetTodosHandler)
			handler.ServeHTTP(recorder, req)

			if status := recorder.Code; status != tc.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

			if todos, ok := tc.expectedBody.([]*models.Todo); ok {
				var decodedTodos []*models.Todo
				if err := json.NewDecoder(recorder.Body).Decode(&decodedTodos); err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				if len(decodedTodos) != len(todos) {
					t.Errorf("Handler returned unexpected body length:\nGot:  %d\nWant: %d", len(decodedTodos), len(todos))
				}

				for i := range todos {
					if *decodedTodos[i] != *todos[i] {
						t.Errorf("Handler returned unexpected body:\nGot:  %+v\nWant: %+v", decodedTodos[i], todos[i])
					}
				}
			} else if errorBody, ok := tc.expectedBody.(map[string]string); ok {
				var decodedErrorBody map[string]string
				if err := json.NewDecoder(recorder.Body).Decode(&decodedErrorBody); err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				if !reflect.DeepEqual(decodedErrorBody, errorBody) {
					t.Errorf("Handler returned unexpected body:\nGot:  %+v\nWant: %+v", decodedErrorBody, errorBody)
				}
			} else {
				t.Fatalf("Unsupported type for expectedBody: %T", tc.expectedBody)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestUpdateTodoHandler(t *testing.T) {
	type testCase struct {
		name             string
		payload          string
		expectedBody     interface{}
		expectedStatus   int
		mockReturn       func(*stores.MockStore)
		token            string
		getUserMockStore func(*stores.MockStore)
		getTodoMockStore func(*stores.MockStore)
	}
	tests := []testCase{
		{
			name:    "Update Todo",
			payload: `{"task_name": "Learn Go", "completed": false, "due_date": "2024-11-30T23:59:59Z"}`,
			expectedBody: &models.Todo{
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
				}).Return(&models.User{}, nil)
			},
			getTodoMockStore: func(mockStore *stores.MockStore) {
				mockStore.On("GetTodo", 1, 0).Return(nil)
			},
		},
		{
			name:           "Invalid Update Todo",
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
				}).Return(&models.User{}, nil)
			},
			getTodoMockStore: func(mockStore *stores.MockStore) {
				mockStore.On("GetTodo", 1, 0).Return(nil)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockStore := stores.InitMockStore()
			tc.mockReturn(mockStore)
			tc.getUserMockStore(mockStore)
			tc.getTodoMockStore(mockStore)
			stores.InitStore(mockStore)

			req, err := http.NewRequest("PUT", "/todos/1", strings.NewReader(tc.payload))
			req.Header.Set("Authorization", "Bearer "+tc.token)
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

			if tc.expectedBody != nil {
				if todo, ok := tc.expectedBody.(*models.Todo); ok {
					var decodedTodo models.Todo
					if err := json.NewDecoder(recorder.Body).Decode(&decodedTodo); err != nil {
						t.Fatalf("Failed to decode response body: %v", err)
					}
					if decodedTodo != *todo {
						t.Errorf("Handler returned unexpected body:\nGot:  %+v\nWant: %+v", decodedTodo, todo)
					}
				} else {
					if errorBody, ok := tc.expectedBody.(map[string]string); ok {
						var decodedErrorBody map[string]string
						if err := json.NewDecoder(recorder.Body).Decode(&decodedErrorBody); err != nil {
							t.Fatalf("Failed to decode response body: %v", err)
						}
						if !reflect.DeepEqual(decodedErrorBody, errorBody) {
							t.Errorf("Handler returned unexpected body:\nGot:  %+v\nWant: %+v", decodedErrorBody, errorBody)
						}
					}
				}
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
