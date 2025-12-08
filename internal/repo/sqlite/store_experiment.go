package sqlite

import (
//	"fmt"
	"database/sql"
//	"log"
//	"github.com/gosimple/slug"
	_ "github.com/glebarez/go-sqlite"

	"go_rest_crud/internal/repo"
	"go_rest_crud/internal/entity"
)


type SQLiteExperimentStore struct {
	*sql.DB
}

// TODO: add the port and it parameters. For now it's only :memory.
func NewSQLiteExperimentStore(db *sql.DB) *SQLiteExperimentStore {
	return &SQLiteExperimentStore{
		db,
	}
}

func (s *SQLiteExperimentStore) Add(e entity.Experiment) error {
	// TODO: Заглушка
	return repo.NotFoundErr
}

func (s *SQLiteExperimentStore) Get(id int) (entity.Experiment, error) {
	// TODO: Заглушка
	e := entity.Experiment{1, "Example Experiment", "Some Description"}
	return e, repo.NotFoundErr
}

func (s *SQLiteExperimentStore) Update(id int, e entity.Experiment) error {
	// TODO: Заглушка
	return repo.NotFoundErr
}

func (s *SQLiteExperimentStore) List() (map[string]entity.Experiment, error) {
    rows, err := s.Query("SELECT id, name, description FROM equipment")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    result := make(map[string]entity.Experiment)
    for rows.Next() {
        var e entity.Experiment
        if err := rows.Scan(&e.ID, &e.Name, &e.Description); err != nil {
            return nil, err
        }
        result[e.Name] = e
    }

    return result, nil
}

func (s *SQLiteExperimentStore) Remove(id int) error {
	// TODO: Заглушка
	return repo.NotFoundErr
}
