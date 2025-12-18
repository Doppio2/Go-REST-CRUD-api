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

type SQLiteEquipmentStore struct {
	*sql.DB
}

// TODO: add the port and it parameters. For now it's only :memory.
func NewSQLiteEquipmentStore(db *sql.DB) *SQLiteEquipmentStore {
	return &SQLiteEquipmentStore{
		db,
	}
}

/*
func (store *SQLiteEquipmentStore) PrintDBSchema() {
	rows, err := store.Query("PRAGMA table_info(equipment);")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dfltValue sql.NullString

		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Column: %s, Type: %s, NotNull: %d, Default: %v, PK: %d\n",
		name, ctype, notnull, dfltValue, pk)
	}
}
*/

func (s *SQLiteEquipmentStore) Add(e entity.Equipment) (int, error) {
    query := "INSERT INTO equipment (name, description, creation_date) VALUES (?, ?, ?)"

	// Используем res для получения LastInsertID.
	res, err := s.Exec(query, e.Name, e.Description, e.CreationDate) 
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil // Возвращаем ID
}

// Стоит ли назвать функцию GetById, а не просто Get???? Не знаю пока. Если других фукнция не планируется, мб и не стоит.
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
	// TODO: нужно ли обновлять дату создания?????
    _, err := s.Exec("UPDATE equipment SET name = ?, description = ? WHERE id = ?", e.Name, e.Description, id)
    return err
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
	// 1. Получаем данные через наш существующий метод List
	equipmentMap, err := s.List()
	if err != nil {
		return fmt.Errorf("ошибка при получении данных: %v", err)
	}

	// 2. Создаем файл
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("не удалось создать файл: %v", err)
	}
	defer file.Close()

	// Добавляем BOM (Byte Order Mark) для корректного отображения кириллицы в Excel
	file.WriteString("\xEF\xBB\xBF")

	writer := csv.NewWriter(file)
	// Устанавливаем точку с запятой как разделитель (стандарт для Excel в RU регионе)
	writer.Comma = ';'
	defer writer.Flush()

	// 3. Записываем "шапку" таблицы
	headers := []string{"ID", "Название", "Описание", "Дата создания"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("ошибка записи заголовков: %v", err)
	}

	// 4. Проходим по мапе и записываем строки
	for _, e := range equipmentMap {
		row := []string{
			strconv.Itoa(e.ID),
			e.Name,
			e.Description,
			e.CreationDate, // Удобный формат даты
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("ошибка записи строки ID %d: %v", e.ID, err)
		}
	}

	return nil
}
