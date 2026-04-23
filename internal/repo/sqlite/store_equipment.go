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

type SQLiteEquipmentStore struct {
	*sql.DB
}

func NewSQLiteEquipmentStore(db *sql.DB) *SQLiteEquipmentStore {
	return &SQLiteEquipmentStore{
		db,
	}
}

func (s *SQLiteEquipmentStore) Add(e entity.Equipment) (int, error) {
	query := "INSERT INTO equipment (name, description, creation_date) VALUES (?, ?, ?)"

	res, err := s.Exec(query, e.Name, e.Description, e.CreationDate)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *SQLiteEquipmentStore) Get(id int) (entity.Equipment, error) {
	var e entity.Equipment
	rows := s.QueryRow("SELECT id, name, description, creation_date FROM equipment WHERE id = ?", id)
	err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.CreationDate)
	if err == sql.ErrNoRows {
		return e, repo.NotFoundErr
	}
	return e, err
}

func (s *SQLiteEquipmentStore) Update(id int, e entity.Equipment) error {
	res, err := s.Exec("UPDATE equipment SET name = ?, description = ? WHERE id = ?", e.Name, e.Description, id)
	if err != nil {
		return err
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		return repo.NotFoundErr
	}

	return nil
}

func (s *SQLiteEquipmentStore) List() (map[int]entity.Equipment, error) {
	rows, err := s.Query("SELECT id, name, description, creation_date FROM equipment")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]entity.Equipment)
	for rows.Next() {
		var e entity.Equipment
		if err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.CreationDate); err != nil {
			return nil, err
		}
		result[e.ID] = e
	}

	return result, nil
}

func (s *SQLiteEquipmentStore) Remove(id int) error {
	res, err := s.Exec("DELETE FROM equipment WHERE id = ?", id)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return repo.NotFoundErr
	}
	return nil
}

func (s *SQLiteEquipmentStore) ExportAllToFile(filePath string) error {
	equipmentMap, err := s.List()
	if err != nil {
		return fmt.Errorf("load equipment for export: %w", err)
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

	headers := []string{"ID", "Название", "Описание", "Дата создания"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("write CSV header: %w", err)
	}

	for _, e := range equipmentMap {
		row := []string{
			strconv.Itoa(e.ID),
			e.Name,
			e.Description,
			e.CreationDate,
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("write CSV row for equipment %d: %w", e.ID, err)
		}
	}

	return nil
}
