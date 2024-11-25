package stores

import (
	"database/sql"
	"testing"
	"time"
	"todo-list/src/models"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type StoreSuite struct {
	suite.Suite
	store *DbStore
	db    *sql.DB
}

func (s *StoreSuite) SetupSuite() {
	connString := "host=localhost port=5432 user=postgres password=secret dbname=todos_test sslmode=disable"
	db, err := sql.Open("postgres", connString)
	if err != nil {
		s.T().Fatal(err)
	}
	s.db = db
	s.store = &DbStore{DB: db}
}

func (s *StoreSuite) SetupTest() {
	_, err := s.db.Query("DELETE FROM todos")
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *StoreSuite) TearDownSuite() {
	s.db.Close()
}

func TestStoreSuite(t *testing.T) {
	s := new(StoreSuite)
	suite.Run(t, s)
}

func (s *StoreSuite) TestCreateTodo() {
	s.store.CreateTodo(&models.Todo{
		TaskName:  "Test task",
		Completed: false,
		DueDate:   time.Now(),
	})

	res, err := s.db.Query(`SELECT COUNT(*) from todos where task_name='Test task' AND completed=false`)
	if err != nil {
		s.T().Fatal(err)
	}

	var count int
	for res.Next() {
		err := res.Scan(&count)
		if err != nil {
			s.T().Fatal(err)
		}
	}

	if count != 1 {
		s.T().Errorf("incorrect count, wanted 1, got %d", count)
	}
}
