package sqlite

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	_ "github.com/glebarez/go-sqlite"

	"go_rest_crud/internal/entity"
	"go_rest_crud/internal/repo"
)

type SQLiteExperimentEquipmentStore struct {
	*sql.DB
}

func NewSQLiteExperimentEquipmentStore(db *sql.DB) *SQLiteExperimentEquipmentStore {
	return &SQLiteExperimentEquipmentStore{
		db,
	}
}

func (s *SQLiteExperimentEquipmentStore) Add(experimentId int, equipmentId int) error {
	query := "INSERT INTO experiment_equipment (experiment_id, equipment_id) VALUES (?, ?)"
	_, err := s.Exec(query, experimentId, equipmentId)
	return err
}

func (s *SQLiteExperimentEquipmentStore) Remove(experimentId int, equipmentId int) error {
	query := "DELETE FROM experiment_equipment WHERE experiment_id = ? AND equipment_id = ?"
	res, err := s.Exec(query, experimentId, equipmentId)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return repo.NotFoundErr
	}

	return nil
}

func (s *SQLiteExperimentEquipmentStore) GetEquipment(experimentId int, equipmentId int) (entity.Equipment, error) {
	var e entity.Equipment

	query := `
	    SELECT equipment.id, equipment.name, equipment.description
		FROM equipment
		JOIN experiment_equipment ON experiment_equipment.equipment_id = equipment.id
		WHERE experiment_equipment.experiment_id = ? AND equipment.id = ?
	`

	row := s.QueryRow(query, experimentId, equipmentId)

	err := row.Scan(&e.ID, &e.Name, &e.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Equipment{}, repo.NotFoundErr
		}
		return entity.Equipment{}, err
	}

	return e, nil
}

func (s *SQLiteExperimentEquipmentStore) ListEquipment(experimentId int) (map[int]entity.Equipment, error) {
	query := `
	    SELECT equipment.id, equipment.name, equipment.description
		FROM equipment
		JOIN experiment_equipment ON experiment_equipment.equipment_id = equipment.id
		WHERE experiment_equipment.experiment_id = ?
	`

	rows, err := s.Query(query, experimentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make(map[int]entity.Equipment)

	for rows.Next() {
		var e entity.Equipment
		if err := rows.Scan(&e.ID, &e.Name, &e.Description); err != nil {
			return nil, err
		}
		res[e.ID] = e
	}

	return res, nil
}

func (s *SQLiteExperimentEquipmentStore) ExportEquipmentToFile(experimentId int, filePath string) error {
	equipmentMap, err := s.ListEquipment(experimentId)
	if err != nil {
		return fmt.Errorf("load experiment equipment for export: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create export file: %w", err)
	}
	defer file.Close()

	file.WriteString("\xEF\xBB\xBF")

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	writer.Write([]string{"Отчет по оборудованию для эксперимента ID:", strconv.Itoa(experimentId)})
	writer.Write([]string{})
	writer.Write([]string{"ID Оборудования", "Название", "Описание"})

	for _, e := range equipmentMap {
		row := []string{
			strconv.Itoa(e.ID),
			e.Name,
			e.Description,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
