package repo

import (
	"fmt"
	"database/sql"
	"log"
	"github.com/gosimple/slug"
	_ "github.com/glebarez/go-sqlite"

	"go_rest_crud/internal/entity"
)

type SQLiteStore struct {
	*sql.DB
}

// TODO: add the port and it parameters. For now it's only :memory.
func NewSQLiteStore(db *sql.DB) *SQLiteStore {
	return &SQLiteStore{
		db,
	}
}

func (store *SQLiteStore) PrintDBSchema() {
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

func (s *SQLiteStore) Add(name string, e entity.Equipment) error {
    query := "INSERT INTO equipment (name, field) VALUES (?, ?)"
    _, err := s.Exec(query, name, e.Field)
    return err
}

func (s *SQLiteStore) Get(name string) (entity.Equipment, error) {
    var e entity.Equipment
    err := s.QueryRow("SELECT id, name, field FROM equipment WHERE name = ?", name).
        Scan(&e.ID, &e.Name, &e.Field)
    if err == sql.ErrNoRows {
        return e, NotFoundErr
    }
    return e, err
}

func (s *SQLiteStore) Update(name string, e entity.Equipment) error {
    newSlug := slug.Make(e.Name)
    _, err := s.Exec("UPDATE equipment SET name = ?, field = ? WHERE name = ?", newSlug, e.Field, name)
    return err
}

func (s *SQLiteStore) List() (map[string]entity.Equipment, error) {
    rows, err := s.Query("SELECT id, name, field FROM equipment")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    result := make(map[string]entity.Equipment)
    for rows.Next() {
        var e entity.Equipment
        if err := rows.Scan(&e.ID, &e.Name, &e.Field); err != nil {
            return nil, err
        }
        result[e.Name] = e
    }
    return result, nil
}

func (s *SQLiteStore) Remove(name string) error {
    res, err := s.Exec("DELETE FROM equipment WHERE name = ?", name)
    if err != nil {
        return err
    }
    count, _ := res.RowsAffected()
    if count == 0 {
        return NotFoundErr
    }
    return nil
}
