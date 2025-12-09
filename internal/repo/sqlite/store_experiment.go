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
    query := "INSERT INTO experiment (name, description) VALUES (?, ?)"
    _, err := s.Exec(query, e.Name, e.Description)
    return err
}

func (s *SQLiteExperimentStore) Get(id int) (entity.Experiment, error) {
    var e entity.Experiment
    err := s.QueryRow("SELECT id, name, description FROM experiment WHERE id = ?", id).
        Scan(&e.ID, &e.Name, &e.Description)
    if err == sql.ErrNoRows {
        return e, repo.NotFoundErr
    }
    return e, err
}

func (s *SQLiteExperimentStore) Update(id int, e entity.Experiment) error {
    _, err := s.Exec("UPDATE experiment SET name = ?, description = ? WHERE id = ?", e.Name, e.Description, id)
    return err
}

func (s *SQLiteExperimentStore) List() (map[int]entity.Experiment, error) {
    rows, err := s.Query("SELECT id, name, description FROM experiment")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    result := make(map[int]entity.Experiment)
    for rows.Next() {
        var e entity.Experiment
        if err := rows.Scan(&e.ID, &e.Name, &e.Description); err != nil {
            return nil, err
        }
        result[e.ID] = e
    }

    return result, nil
}

func (s *SQLiteExperimentStore) Remove(id int) error {
    res, err := s.Exec("DELETE FROM experiment WHERE id = ?", id)
    if err != nil {
        return err
    }
    count, _ := res.RowsAffected()
    if count == 0 {
        return repo.NotFoundErr
    }
    return nil
}
