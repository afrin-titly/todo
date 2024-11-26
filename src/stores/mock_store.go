package stores

import (
	"todo-list/src/models"

	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

// UpdateTodo implements Store.
func (m *MockStore) UpdateTodo(todo *models.Todo, ID int) (*models.Todo, error) {
	panic("unimplemented")
}

func (m *MockStore) CreateTodo(todo *models.Todo) (*models.Todo, error) {
	rets := m.Called(todo)
	return rets.Get(0).(*models.Todo), rets.Error(1)
}

func (m *MockStore) GetTodos() ([]*models.Todo, error) {
	rets := m.Called()
	return rets.Get(0).([]*models.Todo), rets.Error(1)
}

func InitMockStore() *MockStore {
	s := new(MockStore)
	return s
}
