package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"todo-list/src/models"
	"todo-list/src/stores"
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
		name           string
		payload        string
		expectedTodo   *models.Todo
		expectedStatus int
		mockReturn     func(*stores.MockStore)
	}

	fixedTime := time.Date(2024, time.November, 24, 0, 0, 0, 0, time.UTC)

	tests := []testCase{
		{
			name:    "Valid Todo",
			payload: `{"task_name": "Learn Go", "completed": false, "due_date": "2024-11-30T23:59:59Z"}`,
			expectedTodo: &models.Todo{
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
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockStore := stores.InitMockStore()
			tc.mockReturn(mockStore)
			stores.InitStore(mockStore)

			req, err := http.NewRequest("POST", "/todos", strings.NewReader(tc.payload))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(CreateTodoHandler)
			handler.ServeHTTP(recorder, req)

			if status := recorder.Code; status != tc.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

			if tc.expectedTodo != nil {
				var responseTodo models.Todo
				err := json.NewDecoder(recorder.Body).Decode(&responseTodo)
				if err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				if responseTodo != *tc.expectedTodo {
					t.Errorf("Handler returned unexpected body:\nGot:  %+v\nWant: %+v", responseTodo, tc.expectedTodo)
				}
			}

			mockStore.AssertExpectations(t)
		})
	}

}
