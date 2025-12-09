package sqlite

import (
//	"fmt"
	"database/sql"
//	"log"
	_ "github.com/glebarez/go-sqlite"

	"go_rest_crud/internal/repo"
	"go_rest_crud/internal/entity"
)

type SQLiteExperimentEquipmentStore struct {
	*sql.DB
}

// TODO: нужно реализовать методы для взаимодействия со связующей таблицей.

// TODO: add the port and it parameters. For now it's only :memory.
func NewSQLiteExperimentEquipmentStore(db *sql.DB) *SQLiteExperimentEquipmentStore {
	return &SQLiteExperimentEquipmentStore{
		db,
	}
}

func (s *SQLiteExperimentEquipmentStore) Add(experimentId int, equipmentId int) error {
	// TODO: Заглушка
	return repo.NotFoundErr
}

func (s *SQLiteExperimentEquipmentStore) Remove(experimentId int, equipmentId int) error {
	// TODO: Заглушка
	return nil
}

func (s *SQLiteExperimentEquipmentStore) ListEquipment() (map[int]entity.Equipment, error) {
	// TODO: Заглушка
	return nil, nil
}

func (s *SQLiteExperimentEquipmentStore) ListExperiments() (map[int]entity.Experiment, error) {
	// TODO: Заглушка
	return nil, nil
}
