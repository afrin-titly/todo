package stores

import (
	"todo-list/src/models"

	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) CreateTodo(todo *models.Todo, userID int) (*models.Todo, error) {
	rets := m.Called(todo, userID)
	return rets.Get(0).(*models.Todo), rets.Error(1)
}

func (m *MockStore) GetTodos(userID int) ([]*models.Todo, error) {
	rets := m.Called(userID)
	return rets.Get(0).([]*models.Todo), rets.Error(1)
}

func (m *MockStore) GetTodo(todoID int, userID int) error {
	rets := m.Called(todoID, userID)
	return rets.Error(0)
}

func (m *MockStore) UpdateTodo(todo *models.Todo, todoID int) (*models.Todo, error) {
	rets := m.Called(todo, todoID)
	return rets.Get(0).(*models.Todo), rets.Error(1)
}

func (m *MockStore) DeleteTodo(todoID int) error {
	rets := m.Called(todoID)
	return rets.Error(0)
}

func (m *MockStore) CreateUser(user *models.User) (*models.User, error) {
	rets := m.Called(user)
	return rets.Get(0).(*models.User), rets.Error(1)
}

func (m *MockStore) GetUser(user *models.User) (*models.User, error) {
	rets := m.Called(user)
	return rets.Get(0).(*models.User), rets.Error(1)
}

func InitMockStore() *MockStore {
	s := new(MockStore)
	return s
}
