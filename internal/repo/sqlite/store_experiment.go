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

type SQLiteExperimentStore struct {
	*sql.DB
}

func NewSQLiteExperimentStore(db *sql.DB) *SQLiteExperimentStore {
	return &SQLiteExperimentStore{
		db,
	}
}

func (s *SQLiteExperimentStore) Add(ex entity.Experiment) (int, error) {
	query := "INSERT INTO experiment (name, description, creation_date) VALUES (?, ?, ?)"

	res, err := s.Exec(query, ex.Name, ex.Description, ex.CreationDate)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *SQLiteExperimentStore) Get(id int) (entity.Experiment, error) {
	var ex entity.Experiment
	err := s.QueryRow("SELECT id, name, description, creation_date FROM experiment WHERE id = ?", id).
		Scan(&ex.ID, &ex.Name, &ex.Description, &ex.CreationDate)
	if err == sql.ErrNoRows {
		return ex, repo.NotFoundErr
	}
	return ex, err
}

func (s *SQLiteExperimentStore) Update(id int, ex entity.Experiment) error {
	res, err := s.Exec("UPDATE experiment SET name = ?, description = ? WHERE id = ?", ex.Name, ex.Description, id)
	if err != nil {
		return err
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		return repo.NotFoundErr
	}

	return nil
}

func (s *SQLiteExperimentStore) List() (map[int]entity.Experiment, error) {
	rows, err := s.Query("SELECT id, name, description, creation_date FROM experiment")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]entity.Experiment)
	for rows.Next() {
		var ex entity.Experiment
		if err := rows.Scan(&ex.ID, &ex.Name, &ex.Description, &ex.CreationDate); err != nil {
			return nil, err
		}
		result[ex.ID] = ex
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

func (s *SQLiteExperimentStore) ExportAllToFile(filePath string) error {
	experimentsMap, err := s.List()
	if err != nil {
		return fmt.Errorf("load experiments for export: %w", err)
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

	headers := []string{"ID Эксперимента", "Название", "Описание", "Дата проведения"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, ex := range experimentsMap {
		row := []string{
			strconv.Itoa(ex.ID),
			ex.Name,
			ex.Description,
			ex.CreationDate,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
