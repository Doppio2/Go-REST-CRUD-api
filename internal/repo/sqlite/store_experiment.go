package sqlite

import (
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
	// TODO: нужно ли обновлять дату создания?????
    _, err := s.Exec("UPDATE experiment SET name = ?, description = ? WHERE id = ?", ex.Name, ex.Description, id, )
    return err
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
	// TODO: В остальных структурах у меня выражение DELETE FROM ... WHERE в отдельную перменную вынесено. Возможно стоит так же сделать.
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
