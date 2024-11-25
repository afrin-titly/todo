package stores

import (
	"todo-list/src/models"

	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) CreateTodo(todo *models.Todo) (*models.Todo, error) {
	rets := m.Called(todo)
	var todoResult *models.Todo
	if ret := rets.Get(0); ret != nil {
		todoResult = ret.(*models.Todo)
	}
	return todoResult, rets.Error(1)
}

func InitMockStore() *MockStore {
	s := new(MockStore)
	return s
}
