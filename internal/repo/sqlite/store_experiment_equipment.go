package sqlite

import (
	"fmt"
	"os"
	"encoding/csv"
	"strconv"
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
	query := "INSERT INTO experiment_equipment (experiment_id, equipment_id) VALUES (?, ?)"
	// TODO: может быть логирование какое-нибудь?
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

// Метод для получения списка оборудования, которое используется в эксперименте.
func (s *SQLiteExperimentEquipmentStore) ListEquipment(experimentId int) (map[int]entity.Equipment, error) {
	// TODO: позже лучше разобрать.
	query := `
	    SELECT equipment.id, equipment.name, equipment.description
		FROM equipment
		JOIN experiment_equipment ON experiment_equipment.equipment_id = equipment.id  
		WHERE experiment_equipment.experiment_id = ?
	`

	rows, err := s.Query(query, experimentId)
	if err != nil {
		// TODO: позже лучше прологировать.
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

// Метод для получения списка экспериментов, в которых нужно оборудование, id которого мы передаем в метод.
// NOTE: пока что это нигде не используется. И не думаю, что я придумаю, как это использовать.
/*
func (s *SQLiteExperimentEquipmentStore) ListExperiments(equipmentId int) (map[int]entity.Experiment, error) {
	// TODO: позже лучше разобрать.
	query := `
	    SELECT experiment.id, experiment.name, experiment.description
		FROM experiment
		JOIN experiment_equipment ON experiment_equipment.experiment_id = experiment.id  
		WHERE experiment_equipment.equipment_id = ?
	`

	rows, err := s.Query(query, equipmentId)
	if err != nil {
		// TODO: позже лучше прологировать.
		return nil, err
	}
	defer rows.Close()

	res := make(map[int]entity.Experiment)

	for rows.Next() {
		var ex entity.Experiment
		if err := rows.Scan(&ex.ID, ex.Name, ex.Description); err != nil {
			return nil, err
		}
		res[ex.ID] = ex
	}

	return res, nil
}
*/

// Новый метод для экспорта оборудования конкретного эксперимента в файл
func (s *SQLiteExperimentEquipmentStore) ExportEquipmentToFile(experimentId int, filePath string) error {
	// 1. Получаем список оборудования для этого эксперимента
	equipmentMap, err := s.ListEquipment(experimentId)
	if err != nil {
		return fmt.Errorf("ошибка получения данных для эксперимента %d: %w", experimentId, err)
	}

	// 2. Создаем файл
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer file.Close()

	// Добавляем BOM для кириллицы
	file.WriteString("\xEF\xBB\xBF")

	writer := csv.NewWriter(file)
	writer.Comma = ';' // Для Excel
	defer writer.Flush()

	// 3. Записываем информацию о том, чей это отчет
	writer.Write([]string{"Отчет по оборудованию для эксперимента ID:", strconv.Itoa(experimentId)})
	writer.Write([]string{}) // Пустая строка для красоты

	// 4. Заголовки таблицы
	writer.Write([]string{"ID Оборудования", "Название", "Описание"})

	// 5. Заполняем данными
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
