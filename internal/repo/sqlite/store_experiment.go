package sqlite

import (
	"fmt"
	"os"
	"encoding/csv"
	"strconv"
	"database/sql"

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

func (s *SQLiteExperimentStore) Add(ex entity.Experiment) (int, error) {
    query := "INSERT INTO experiment (name, description, creation_date) VALUES (?, ?, ?)"

	// Используем res для получения LastInsertID.
	res, err := s.Exec(query, ex.Name, ex.Description, ex.CreationDate) 
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil // Возвращаем ID
}

// Метод для получения данных об эксперименте по ID.
func (s *SQLiteExperimentStore) Get(id int) (entity.Experiment, error) {
    var ex entity.Experiment
    err := s.QueryRow("SELECT id, name, description, creation_date FROM experiment WHERE id = ?", id).
        Scan(&ex.ID, &ex.Name, &ex.Description, &ex.CreationDate)
    if err == sql.ErrNoRows {
        return ex, repo.NotFoundErr
    }
    return ex, err
}

// Метод для обновления данных об эксперименте.
func (s *SQLiteExperimentStore) Update(id int, ex entity.Experiment) error {
    _, err := s.Exec("UPDATE experiment SET name = ?, description = ? WHERE id = ?", ex.Name, ex.Description, id, )
    return err
}

// Метод для получения списках всех экспериментов из базы данных.
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

// Метод для удаления данных об эксперименте по ID.
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

// ExportAllToFile для экспериментов
func (s *SQLiteExperimentStore) ExportAllToFile(filePath string) error {
	// 1. Получаем данные
	experimentsMap, err := s.List()
	if err != nil {
		return fmt.Errorf("ошибка при получении списка экспериментов: %w", err)
	}

	// 2. Создаем файл
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("не удалось создать файл: %w", err)
	}
	defer file.Close()

	// Добавляем BOM для корректного отображения кириллицы в Excel
	file.WriteString("\xEF\xBB\xBF")

	writer := csv.NewWriter(file)
	writer.Comma = ';' // Точка с запятой для корректного открытия в Excel (RU)
	defer writer.Flush()

	// 3. Заголовки таблицы
	headers := []string{"ID Эксперимента", "Название", "Описание", "Дата проведения"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// 4. Заполнение данными
	for _, ex := range experimentsMap {
		row := []string{
			strconv.Itoa(ex.ID),
			ex.Name,
			ex.Description,
			ex.CreationDate, // Обычно для экспериментов достаточно даты без времени
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
