package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"todo-list/src/models"
	"todo-list/src/stores"
)

func TestCreateUserHandler(t *testing.T) {
	type testCase struct {
		name           string
		payload        string
		expectedUser   *models.User
		expectedStatus int
		mockStore      func(*stores.MockStore)
	}

	tests := []testCase{
		{
			name:    "Create user successfully",
			payload: `{"email": "test@mail.com", "password": "password"}`,
			expectedUser: &models.User{
				UserName: "testuser",
				Email:    "test@mail.com",
				Password: "password",
			},
			expectedStatus: http.StatusCreated,
			mockStore: func(mockStore *stores.MockStore) {
				mockStore.On("CreateUser", &models.User{
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
			name:           "Invalid Email/Password",
			payload:        `{"email": "testmail.com", "password": "pass"}`,
			expectedUser:   &models.User{},
			expectedStatus: http.StatusBadRequest,
			mockStore:      func(mockStore *stores.MockStore) {},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockStore := stores.InitMockStore()
			tc.mockStore(mockStore)
			stores.InitStore(mockStore)

			req, err := http.NewRequest("POST", "/users", strings.NewReader(tc.payload))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(CreateUserHandler)
			handler.ServeHTTP(recorder, req)
			t.Logf("-------Handler returned body-------: %v", recorder.Body.String())

			if status := recorder.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

			if tc.expectedStatus == http.StatusOK {
				var user models.User
				err := json.NewDecoder(recorder.Body).Decode(&user)

				if err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				if !reflect.DeepEqual(&user, tc.expectedUser) {
					t.Errorf("handler returned unexpected body: got %v want %v", user, tc.expectedUser)
				}
			}
			mockStore.AssertExpectations(t)
		})
	}
}

func TestLoginUserHandler(t *testing.T) {
	type testCase struct {
		name           string
		payload        string
		expectedStatus int
		mockStore      func(*stores.MockStore)
	}

	tests := []testCase{
		{
			name:           "Login user successfully",
			payload:        `{"email": "test@mail.com", "password": "password"}`,
			expectedStatus: http.StatusOK,
			mockStore: func(mockStore *stores.MockStore) {
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
			name:           "Unsuccessful login",
			payload:        `{"email": "nobody@mail.com", "password": "password"}`,
			expectedStatus: http.StatusUnauthorized,
			mockStore: func(mockStore *stores.MockStore) {
				mockStore.On("GetUser", &models.User{
					Email:    "nobody@mail.com",
					Password: "password",
				}).Return(&models.User{}, sql.ErrNoRows)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockStore := stores.InitMockStore()
			tc.mockStore(mockStore)
			stores.InitStore(mockStore)

			req, err := http.NewRequest("POST", "/users/login", strings.NewReader(tc.payload))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(LoginUserHandler)
			handler.ServeHTTP(recorder, req)
			t.Logf("-------Handler returned body-------: %v", recorder.Body.String())

			if status := recorder.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}
			mockStore.AssertExpectations(t)
		})
	}
}
