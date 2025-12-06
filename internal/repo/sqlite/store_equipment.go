package sqlite

import (
	"fmt"
	"database/sql"
	"log"
	"github.com/gosimple/slug"
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

func (s *SQLiteEquipmentStore) Add(name string, e entity.Equipment) error {
    query := "INSERT INTO equipment (name, description) VALUES (?, ?)"
    _, err := s.Exec(query, name, e.Description)
    return err
}

func (s *SQLiteEquipmentStore) Get(name string) (entity.Equipment, error) {
    var e entity.Equipment
    err := s.QueryRow("SELECT id, name, description FROM equipment WHERE name = ?", name).
        Scan(&e.ID, &e.Name, &e.Description)
    if err == sql.ErrNoRows {
        return e, repo.NotFoundErr
    }
    return e, err
}

func (s *SQLiteEquipmentStore) Update(name string, e entity.Equipment) error {
    newSlug := slug.Make(e.Name)
    _, err := s.Exec("UPDATE equipment SET name = ?, description = ? WHERE name = ?", newSlug, e.Description, name)
    return err
}

func (s *SQLiteEquipmentStore) List() (map[string]entity.Equipment, error) {
    rows, err := s.Query("SELECT id, name, description FROM equipment")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    result := make(map[string]entity.Equipment)
    for rows.Next() {
        var e entity.Equipment
        if err := rows.Scan(&e.ID, &e.Name, &e.Description); err != nil {
            return nil, err
        }
        result[e.Name] = e
    }
    return result, nil
}

func (s *SQLiteEquipmentStore) Remove(name string) error {
    res, err := s.Exec("DELETE FROM equipment WHERE name = ?", name)
    if err != nil {
        return err
    }
    count, _ := res.RowsAffected()
    if count == 0 {
        return repo.NotFoundErr
    }
    return nil
}
