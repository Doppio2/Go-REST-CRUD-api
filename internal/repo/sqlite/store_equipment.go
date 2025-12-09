package sqlite

import (
	"fmt"
	"database/sql"
	"log"
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

func (s *SQLiteEquipmentStore) Add(e entity.Equipment) error {
    query := "INSERT INTO equipment (name, description) VALUES (?, ?)"
    _, err := s.Exec(query, e.Name, e.Description)
    return err
}

// Стоит ли назвать функцию GetById, а не просто Get???? Не знаю пока. Если других фукнция не планируется, мб и не стоит.
func (s *SQLiteEquipmentStore) Get(id int) (entity.Equipment, error) {
    var e entity.Equipment
    err := s.QueryRow("SELECT id, name, description FROM equipment WHERE id = ?", id).
        Scan(&e.ID, &e.Name, &e.Description)
    if err == sql.ErrNoRows {
        return e, repo.NotFoundErr
    }
    return e, err
}

func (s *SQLiteEquipmentStore) Update(id int, e entity.Equipment) error {
    _, err := s.Exec("UPDATE equipment SET name = ?, description = ? WHERE id = ?", e.Name, e.Description, id)
    return err
}

func (s *SQLiteEquipmentStore) List() (map[int]entity.Equipment, error) {
    rows, err := s.Query("SELECT id, name, description FROM equipment")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    result := make(map[int]entity.Equipment)
    for rows.Next() {
        var e entity.Equipment
        if err := rows.Scan(&e.ID, &e.Name, &e.Description); err != nil {
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
