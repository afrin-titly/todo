package stores

import (
	"todo-list/src/models"

	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) CreateTodo(todo *models.Todo, ID int) (*models.Todo, error) {
	rets := m.Called(todo)
	return rets.Get(0).(*models.Todo), rets.Error(1)
}

func (m *MockStore) GetTodos() ([]*models.Todo, error) {
	rets := m.Called()
	return rets.Get(0).([]*models.Todo), rets.Error(1)
}

func (m *MockStore) UpdateTodo(todo *models.Todo, ID int) (*models.Todo, error) {
	rets := m.Called(todo, ID)
	return rets.Get(0).(*models.Todo), rets.Error(1)
}

func (m *MockStore) DeleteTodo(ID int) error {
	rets := m.Called(ID)
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
