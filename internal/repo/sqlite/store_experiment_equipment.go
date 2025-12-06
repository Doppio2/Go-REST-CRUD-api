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

type SQLiteExperimentEquipmentStore struct {
	*sql.DB
}

// TODO: add the port and it parameters. For now it's only :memory.
func NewSQLiteExperimentEquipmentStore(db *sql.DB) *SQLiteExperimentEquipmentStore {
	return &SQLiteExperimentEquipmentStore{
		db,
	}
}

func (s *SQLiteExperimentEquipmentStore) Add(name string, e entity.ExperimentEquipment) error {
	// TODO: Заглушка
	return repo.NotFoundErr
}

func (s *SQLiteExperimentEquipmentStore) Get(name string) (entity.ExperimentEquipment, error) {
	// TODO: Заглушка
	e := entity.ExperimentEquipment{1, 2}
	return e, repo.NotFoundErr
}

func (s *SQLiteExperimentEquipmentStore) Update(name string, e entity.ExperimentEquipment) error {
	// TODO: Заглушка
	return repo.NotFoundErr
}

func (s *SQLiteExperimentEquipmentStore) List() (map[string]entity.ExperimentEquipment, error) {
	// TODO: Заглушка
    rows, err := s.Query("SELECT id, name, description FROM equipment")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    result := make(map[string]entity.ExperimentEquipment)
    for rows.Next() {
        var e entity.ExperimentEquipment
        if err := rows.Scan(&e.EquipmentID, &e.ExperimentID); err != nil {
            return nil, err
        }
    }

    return result, nil
}

func (s *SQLiteExperimentEquipmentStore) Remove(name string) error {
	// TODO: Заглушка
	return repo.NotFoundErr
}
